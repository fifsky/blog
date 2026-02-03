package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"app/config"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GuestbookEntry represents a guestbook entry to import
type GuestbookEntry struct {
	Name      string
	Content   string
	IP        string
	CreatedAt time.Time
}

// CDXRecord represents a record from Wayback Machine CDX API
type CDXRecord struct {
	Timestamp   string
	OriginalURL string
	StatusCode  string
}

func main() {
	ctx := context.Background()

	// Initialize config and database
	conf := config.New()
	db, cleanup := config.NewBlogDB(conf)
	defer cleanup()

	log.Println("Starting guestbook import from Wayback Machine...")

	// Fetch archived URLs from CDX API
	records, err := fetchCDXRecords()
	if err != nil {
		log.Fatalf("Failed to fetch CDX records: %v", err)
	}
	log.Printf("Found %d archived records", len(records))

	// Filter valid records (status 200 only)
	validRecords := filterValidRecords(records)
	log.Printf("Filtered to %d valid records with status 200", len(validRecords))

	// Use map for deduplication: key = name+content (skip if same nickname and content)
	entries := make(map[string]GuestbookEntry)

	// Process each archived page
	for i, record := range validRecords {
		url := buildWaybackURL(record.Timestamp, record.OriginalURL)
		log.Printf("[%d/%d] Processing: %s", i+1, len(validRecords), url)

		pageEntries, err := extractEntriesFromPage(url, record.OriginalURL)
		if err != nil {
			log.Printf("Error processing %s: %v", url, err)
			continue
		}

		log.Printf("  Found %d entries on this page", len(pageEntries))

		for _, entry := range pageEntries {
			// Deduplicate by name + content only
			key := fmt.Sprintf("%s|%s", entry.Name, entry.Content)
			if _, exists := entries[key]; !exists {
				entries[key] = entry
			} else {
				log.Printf("    [跳过] 重复内容: %s - %s", entry.Name, truncateString(entry.Content, 30))
			}
		}

		// Rate limiting: 1 second between requests
		time.Sleep(time.Second)
	}

	log.Printf("Extracted %d unique entries", len(entries))

	// Insert entries into database
	inserted := 0
	for _, entry := range entries {
		if err := insertEntry(ctx, db, entry); err != nil {
			log.Printf("Failed to insert entry: %v", err)
			continue
		}
		inserted++
	}

	log.Printf("Successfully inserted %d entries into guestbook table", inserted)
}

// fetchCDXRecords fetches archived URLs from Wayback Machine CDX API
func fetchCDXRecords() ([]CDXRecord, error) {
	url := "https://web.archive.org/cdx/search/cdx?url=http://www.windiness.com/guestbook/index.php*&output=json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var data [][]string
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Skip header row
	if len(data) < 2 {
		return nil, fmt.Errorf("no records found")
	}

	records := make([]CDXRecord, 0, len(data)-1)
	for i := 1; i < len(data); i++ {
		row := data[i]
		if len(row) >= 5 {
			records = append(records, CDXRecord{
				Timestamp:   row[1],
				OriginalURL: row[2],
				StatusCode:  row[4],
			})
		}
	}

	return records, nil
}

// filterValidRecords filters records to only include those with status 200
func filterValidRecords(records []CDXRecord) []CDXRecord {
	valid := make([]CDXRecord, 0)
	seen := make(map[string]bool)

	for _, r := range records {
		if r.StatusCode == "200" {
			// Deduplicate by timestamp+url
			key := r.Timestamp + r.OriginalURL
			if !seen[key] {
				seen[key] = true
				valid = append(valid, r)
			}
		}
	}
	return valid
}

// buildWaybackURL constructs the full Wayback Machine URL
func buildWaybackURL(timestamp, originalURL string) string {
	return fmt.Sprintf("https://web.archive.org/web/%s/%s", timestamp, originalURL)
}

// hasMODParameter checks if URL contains MOD parameter
func hasMODParameter(originalURL string) bool {
	return strings.Contains(strings.ToUpper(originalURL), "MOD=")
}

// extractEntriesFromPage downloads and parses a page to extract guestbook entries
func extractEntriesFromPage(url, originalURL string) ([]GuestbookEntry, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Convert GBK to UTF-8
	html, err := gbkToUtf8(body)
	if err != nil {
		// If conversion fails, try using the original bytes as UTF-8
		log.Printf("  Warning: GBK conversion failed, using original encoding: %v", err)
		html = string(body)
	}

	// Choose extraction method based on URL pattern
	if hasMODParameter(originalURL) {
		return parseWithMODFormat(html)
	}
	return parseLegacyFormat(html)
}

// gbkToUtf8 converts GBK encoded bytes to UTF-8 string
func gbkToUtf8(data []byte) (string, error) {
	reader := transform.NewReader(strings.NewReader(string(data)), simplifiedchinese.GBK.NewDecoder())
	result, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// parseWithMODFormat extracts entries from pages with MOD parameter (2010 style)
// 昵称从 <img src="...oicq.gif" alt="五月的雪 的 QQ 号码：448631595"> 提取
// 日期从 <font style="...">Time: 2010-06-04 22:10:29</font> 提取
// 正文从 <td width="421">content</td> 提取
func parseWithMODFormat(html string) ([]GuestbookEntry, error) {
	entries := make([]GuestbookEntry, 0)

	// Extract nicknames from oicq.gif image alt attribute
	// Pattern: alt="昵称 的 QQ 号码：..." or alt="昵称 的 ..."
	nicknamePattern := regexp.MustCompile(`<img[^>]*oicq\.gif[^>]*alt="([^"]*?) 的 [^"]*"[^>]*>`)
	nicknames := nicknamePattern.FindAllStringSubmatch(html, -1)

	// Extract times: Time: 2010-06-04 22:10:29
	timePattern := regexp.MustCompile(`Time:\s*(\d{4}-\d{1,2}-\d{1,2}\s+\d{1,2}:\d{2}:\d{2})`)
	times := timePattern.FindAllStringSubmatch(html, -1)

	// Extract contents: from <td width="421" ...>content</td>
	// Using [\s\S]*? to match content including newlines
	contentPattern := regexp.MustCompile(`<td[^>]*width\s*=\s*["']?421["']?[^>]*>([\s\S]*?)</td>`)
	contents := contentPattern.FindAllStringSubmatch(html, -1)

	log.Printf("  MOD format: found %d nicknames, %d times, %d contents", len(nicknames), len(times), len(contents))

	// Match entries by position (same index)
	minLen := len(nicknames)
	if len(times) < minLen {
		minLen = len(times)
	}
	if len(contents) < minLen {
		minLen = len(contents)
	}

	for i := 0; i < minLen; i++ {
		nickname := strings.TrimSpace(nicknames[i][1])
		timeStr := strings.TrimSpace(times[i][1])
		content := cleanContent(contents[i][1])

		createdAt, err := parseTime(timeStr)
		if err != nil {
			log.Printf("Failed to parse time '%s': %v", timeStr, err)
			continue
		}

		if nickname == "" {
			continue
		}

		entries = append(entries, GuestbookEntry{
			Name:      nickname,
			Content:   content,
			IP:        "", // IP not easily extractable in this format
			CreatedAt: createdAt,
		})

		// Debug output
		log.Printf("    [%d] 昵称: %s", i+1, nickname)
		log.Printf("        时间: %s", timeStr)
		log.Printf("        内容: %s", truncateString(content, 100))
	}

	return entries, nil
}

// parseLegacyFormat extracts entries from pages without MOD parameter (2007 style)
// 昵称从 <span class="name">皇家澜澜</span> 提取
// 日期从 <span class="input_time">2007-10-12 09:48:38</span> 提取
// 正文从 <div class="content" style="...">...</div> 提取
func parseLegacyFormat(html string) ([]GuestbookEntry, error) {
	entries := make([]GuestbookEntry, 0)

	// Extract nicknames: <span class="name">昵称</span>
	nicknamePattern := regexp.MustCompile(`<span\s+class="name">([^<]*)</span>`)
	nicknames := nicknamePattern.FindAllStringSubmatch(html, -1)

	// Extract times: <span class="input_time">2007-05-11 21:31:08</span>
	timePattern := regexp.MustCompile(`<span\s+class="input_time">([^<]*)</span>`)
	times := timePattern.FindAllStringSubmatch(html, -1)

	// Extract contents: <div class="content" ...>内容</div>
	contentPattern := regexp.MustCompile(`<div\s+class="content"[^>]*>([\s\S]*?)</div>`)
	contents := contentPattern.FindAllStringSubmatch(html, -1)

	log.Printf("  Legacy format: found %d nicknames, %d times, %d contents", len(nicknames), len(times), len(contents))

	// Match entries by position (same index)
	minLen := len(nicknames)
	if len(times) < minLen {
		minLen = len(times)
	}
	if len(contents) < minLen {
		minLen = len(contents)
	}

	for i := 0; i < minLen; i++ {
		nickname := strings.TrimSpace(nicknames[i][1])
		timeStr := strings.TrimSpace(times[i][1])
		content := cleanContent(contents[i][1])

		createdAt, err := parseTime(timeStr)
		if err != nil {
			log.Printf("Failed to parse time '%s': %v", timeStr, err)
			continue
		}

		if nickname == "" {
			continue
		}

		entries = append(entries, GuestbookEntry{
			Name:      nickname,
			Content:   content,
			IP:        "", // IP not easily extractable in this format
			CreatedAt: createdAt,
		})

		// Debug output
		log.Printf("    [%d] 昵称: %s", i+1, nickname)
		log.Printf("        时间: %s", timeStr)
		log.Printf("        内容: %s", truncateString(content, 100))
	}

	return entries, nil
}

// parseTime parses various time formats
func parseTime(timeStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-1-2 15:04:05",
		"2006-01-2 15:04:05",
		"2006-1-02 15:04:05",
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

// cleanContent removes HTML tags and cleans up whitespace
func cleanContent(s string) string {
	// Remove HTML tags
	tagPattern := regexp.MustCompile(`<[^>]+>`)
	s = tagPattern.ReplaceAllString(s, " ")

	// Decode common HTML entities
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&#160;", " ")

	// Clean up whitespace
	spacePattern := regexp.MustCompile(`\s+`)
	s = spacePattern.ReplaceAllString(s, " ")

	return strings.TrimSpace(s)
}

// insertEntry inserts a guestbook entry into the database
func insertEntry(ctx context.Context, db *sql.DB, entry GuestbookEntry) error {
	_, err := db.ExecContext(ctx,
		"INSERT INTO guestbook (name, content, ip, created_at) VALUES (?, ?, ?, ?)",
		entry.Name, entry.Content, entry.IP, entry.CreatedAt,
	)
	return err
}

// truncateString truncates a string to maxLen characters and adds "..." if truncated
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
