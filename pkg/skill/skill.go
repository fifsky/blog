package skill

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

type Frontmatter struct {
	Name          string         `yaml:"name"`
	Description   string         `yaml:"description"`
	License       string         `yaml:"license"`
	Compatibility string         `yaml:"compatibility"`
	AllowedTools  string         `yaml:"allowed_tools"`
	Metadata      map[string]any `yaml:"metadata"`
}

type Resources struct {
	References map[string]string
	Assets     map[string][]byte
	Scripts    map[string]string
}

type Skill struct {
	Frontmatter Frontmatter
	Content     string
	Path        string
	Resources   Resources
}

func (s *Skill) Name() string {
	return s.Frontmatter.Name
}

func (s *Skill) Description() string {
	return s.Frontmatter.Description
}

type Manager struct {
	skillsPath string
	skills     map[string]*Skill
}

func NewManager(skillsPath string) *Manager {
	return &Manager{
		skillsPath: skillsPath,
		skills:     make(map[string]*Skill),
	}
}

func (m *Manager) Load() error {
	m.skills = make(map[string]*Skill)

	entries, err := os.ReadDir(m.skillsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillDir := filepath.Join(m.skillsPath, entry.Name())
		skill, err := ParseSkillDir(skillDir)
		if err != nil {
			continue
		}

		if skill.Frontmatter.Name == "" {
			skill.Frontmatter.Name = entry.Name()
		}

		m.skills[skill.Frontmatter.Name] = skill
	}

	return nil
}

func ParseSkillDir(dir string) (*Skill, error) {
	skillFile := filepath.Join(dir, "SKILL.md")
	if _, err := os.Stat(skillFile); os.IsNotExist(err) {
		skillFile = filepath.Join(dir, "skill.md")
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("SKILL.md not found in %s", dir)
		}
	}

	skill, err := ParseSkillFile(skillFile)
	if err != nil {
		return nil, err
	}

	skill.Path = dir

	// Load resources
	fsys := os.DirFS(dir)
	skill.Resources.References, _ = loadDirFiles(fsys, "references")
	skill.Resources.Scripts, _ = loadDirFiles(fsys, "scripts")
	skill.Resources.Assets, _ = loadDirBinaryFiles(fsys, "assets")

	return skill, nil
}

func loadDirFiles(fsys fs.FS, dir string) (map[string]string, error) {
	files := make(map[string]string)
	err := fs.WalkDir(fsys, dir, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return fs.SkipDir
			}
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(filePath, dir+"/")
		files[rel] = string(b)
		return nil
	})
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	return files, nil
}

func loadDirBinaryFiles(fsys fs.FS, dir string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	err := fs.WalkDir(fsys, dir, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return fs.SkipDir
			}
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(filePath, dir+"/")
		files[rel] = b
		return nil
	})
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	return files, nil
}

func ParseSkillFile(path string) (*Skill, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(contentBytes)
	skill := &Skill{Path: path}

	firstLine, rest, hasMore := cutLine(content)
	if !isFrontmatterDelimiterLine(firstLine) {
		skill.Content = strings.TrimSpace(content)
		return skill, nil
	}
	if !hasMore {
		skill.Content = strings.TrimSpace(content)
		return skill, nil
	}

	search := rest
	frontmatterLen := 0
	foundEnd := false
	for {
		line, remaining, hasNext := cutLine(search)
		if isFrontmatterDelimiterLine(line) {
			search = remaining
			foundEnd = true
			break
		}
		if !hasNext {
			break
		}
		frontmatterLen += len(line) + 1
		search = remaining
	}

	if foundEnd {
		frontmatterContent := rest[:frontmatterLen]
		if err := yaml.Unmarshal([]byte(frontmatterContent), &skill.Frontmatter); err != nil {
			return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
		}
		skill.Content = strings.TrimSpace(search)
	} else {
		skill.Content = strings.TrimSpace(content)
	}

	return skill, nil
}

func cutLine(content string) (line string, rest string, hasMore bool) {
	idx := strings.IndexByte(content, '\n')
	if idx < 0 {
		return content, "", false
	}
	return content[:idx], content[idx+1:], true
}

func isFrontmatterDelimiterLine(line string) bool {
	return trimTrailingCarriageReturn(line) == "---"
}

func trimTrailingCarriageReturn(line string) string {
	if strings.HasSuffix(line, "\r") {
		return line[:len(line)-1]
	}
	return line
}

func (m *Manager) GetSkills() map[string]*Skill {
	return m.skills
}

func (m *Manager) GetSkill(name string) (*Skill, bool) {
	s, ok := m.skills[name]
	return s, ok
}

func (m *Manager) ExecuteScript(ctx context.Context, skillName, scriptPath string, args []string, env map[string]string, timeoutSeconds int) string {
	skill, ok := m.GetSkill(skillName)
	if !ok {
		return mustJSON(map[string]any{
			"error":      fmt.Sprintf("Skill %q not found", skillName),
			"error_code": "SKILL_NOT_FOUND",
		})
	}

	_, ok = skill.Resources.Scripts[scriptPath]
	if !ok {
		return mustJSON(map[string]any{
			"error":      fmt.Sprintf("Script %q not found in skill %q", scriptPath, skillName),
			"error_code": "SCRIPT_NOT_FOUND",
		})
	}

	tmpRoot, err := os.MkdirTemp("", "blades-skill-*")
	if err != nil {
		return mustJSON(map[string]any{
			"error":      fmt.Sprintf("Failed to prepare skill workspace: %v", err),
			"error_code": "WORKSPACE_ERROR",
		})
	}
	defer os.RemoveAll(tmpRoot)

	if err := materializeSkillWorkspace(tmpRoot, skill.Resources); err != nil {
		return mustJSON(map[string]any{
			"error":      fmt.Sprintf("Failed to materialize skill workspace: %v", err),
			"error_code": "WORKSPACE_ERROR",
		})
	}

	fullScriptPath := filepath.Join(tmpRoot, "scripts", scriptPath)
	return executeSkillScript(ctx, tmpRoot, skillName, fullScriptPath, args, env, timeoutSeconds)
}

func materializeSkillWorkspace(root string, resources Resources) error {
	for rel, content := range resources.References {
		if err := writeWorkspaceFile(root, "references", rel, []byte(content), 0o644); err != nil {
			return err
		}
	}
	for rel, content := range resources.Assets {
		if err := writeWorkspaceFile(root, "assets", rel, content, 0o644); err != nil {
			return err
		}
	}
	for rel, content := range resources.Scripts {
		if err := writeWorkspaceFile(root, "scripts", rel, []byte(content), 0o755); err != nil {
			return err
		}
	}
	return nil
}

func writeWorkspaceFile(root string, dir string, rel string, content []byte, mode fs.FileMode) error {
	fullPath := filepath.Join(root, dir, rel)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, content, mode)
}

func executeSkillScript(ctx context.Context, tmpRoot string, skillName string, scriptPath string, args []string, env map[string]string, timeoutSeconds int) string {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 300
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	commandName := scriptPath
	commandArgs := append([]string{}, args...)
	switch strings.ToLower(path.Ext(scriptPath)) {
	case ".py":
		commandName = "python3"
		commandArgs = append([]string{scriptPath}, commandArgs...)
	case ".sh", ".bash":
		commandName = "bash"
		commandArgs = append([]string{scriptPath}, commandArgs...)
	}

	cmd := exec.CommandContext(timeoutCtx, commandName, commandArgs...)
	cmd.Dir = tmpRoot

	cmdEnv := os.Environ()
	for k, v := range env {
		cmdEnv = append(cmdEnv, k+"="+v)
	}
	cmd.Env = cmdEnv

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	status := "success"
	if err != nil {
		switch {
		case errors.Is(timeoutCtx.Err(), context.DeadlineExceeded):
			exitCode = -1
			status = "timeout"
		default:
			status = "error"
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				exitCode = exitErr.ExitCode()
			} else {
				return mustJSON(map[string]any{
					"error":      fmt.Sprintf("Failed to execute script %q: %v", scriptPath, err),
					"error_code": "EXECUTION_ERROR",
				})
			}
		}
	}

	return mustJSON(map[string]any{
		"skill_name":  skillName,
		"script_path": scriptPath,
		"args":        args,
		"stdout":      stdout.String(),
		"stderr":      stderr.String(),
		"exit_code":   exitCode,
		"status":      status,
	})
}

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
