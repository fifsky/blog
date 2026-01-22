package store

import (
	"fmt"
	"strings"
)

// In 是一个泛型辅助函数，用于生成 SQL IN 子句所需的占位符和参数。
// 它接收一个任意类型的切片和起始索引，返回以逗号分隔的占位符字符串和对应的参数切片。
// startIndex 用于指定占位符的起始编号，例如 startIndex=1 生成 $1,$2,$3
// 示例:
//
//	ids := []int{1, 2, 3}
//	placeholders, args := In(ids, 1)
//	// placeholders = "$1,$2,$3"
//	// args = []any{1, 2, 3}
//	query := "SELECT * FROM users WHERE id IN (" + placeholders + ")"
//	db.Query(query, args...)
func In[T any](s []T, startIndex int) (string, []any) {
	if len(s) == 0 {
		return "", nil
	}

	args := make([]any, len(s))
	for i := range s {
		args[i] = s[i]
	}
	return placeholders(len(s), startIndex), args
}

func placeholder(n int) string {
	return fmt.Sprintf("$%d", n)
}

func placeholders(n int, startIndex int) string {
	list := make([]string, 0, n)
	for i := 0; i < n; i++ {
		list = append(list, placeholder(startIndex+i))
	}
	return strings.Join(list, ", ")
}
