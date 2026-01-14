package store

import (
	"strings"
)

// In 是一个泛型辅助函数，用于生成 SQL IN 子句所需的占位符和参数。
// 它接收一个任意类型的切片，返回以逗号分隔的占位符字符串和对应的参数切片。
// 示例:
//
//	ids := []int{1, 2, 3}
//	placeholders, args := In(ids)
//	// placeholders = "?,?,?"
//	// args = []any{1, 2, 3}
//	query := "SELECT * FROM users WHERE id IN (" + placeholders + ")"
//	db.Query(query, args...)
func In[T any](s []T) (string, []any) {
	if len(s) == 0 {
		return "", nil
	}

	placeholders := make([]string, len(s))
	args := make([]any, len(s))
	for i := range placeholders {
		placeholders[i] = "?"
		args[i] = s[i]
	}
	return strings.Join(placeholders, ","), args
}
