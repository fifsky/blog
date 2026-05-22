package aiagent

import (
	"sync"
	"time"

	"github.com/openai/openai-go/v3"
)

type memoryEntry struct {
	messages  []openai.ChatCompletionMessageParamUnion
	expiresAt time.Time
}

// Memory 保存带过期时间的对话历史。
type Memory struct {
	ttl         time.Duration
	maxMessages int
	now         func() time.Time
	cache       sync.Map
}

// NewMemory 创建对话历史缓存。
func NewMemory(ttl time.Duration, maxMessages int) *Memory {
	return &Memory{
		ttl:         ttl,
		maxMessages: maxMessages,
		now:         time.Now,
	}
}

// Get 返回指定 key 的历史消息副本。
func (m *Memory) Get(key string) []openai.ChatCompletionMessageParamUnion {
	if key == "" {
		return nil
	}

	value, ok := m.cache.Load(key)
	if !ok {
		return nil
	}

	entry := value.(*memoryEntry)
	if m.now().After(entry.expiresAt) {
		m.cache.Delete(key)
		return nil
	}

	messages := make([]openai.ChatCompletionMessageParamUnion, len(entry.messages))
	copy(messages, entry.messages)
	return messages
}

// Save 保存指定 key 的历史消息，并按最大数量裁剪。
func (m *Memory) Save(key string, messages []openai.ChatCompletionMessageParamUnion) {
	if key == "" {
		return
	}

	newMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	copy(newMessages, messages)
	if m.maxMessages > 0 && len(newMessages) > m.maxMessages {
		newMessages = newMessages[len(newMessages)-m.maxMessages:]
	}

	m.cache.Store(key, &memoryEntry{
		messages:  newMessages,
		expiresAt: m.now().Add(m.ttl),
	})
}

// Append 追加消息到指定 key 的历史中。
func (m *Memory) Append(key string, messages ...openai.ChatCompletionMessageParamUnion) {
	existing := m.Get(key)
	existing = append(existing, messages...)
	m.Save(key, existing)
}

// CleanExpired 删除已过期的历史消息。
func (m *Memory) CleanExpired() {
	now := m.now()
	m.cache.Range(func(key, value any) bool {
		entry := value.(*memoryEntry)
		if now.After(entry.expiresAt) {
			m.cache.Delete(key)
		}
		return true
	})
}
