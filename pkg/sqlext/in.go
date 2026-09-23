package sqlext

import (
	"strings"
)

// In 生成 SQL IN 子句所需的占位符和参数，用于 Builder 无法生成的非 SELECT 语句
// （例如 DELETE、批量 UPDATE）。它接收一个任意类型的切片，返回以逗号分隔的占位符
// 字符串和对应的参数切片：
//
//	placeholders, args := sqlext.In([]int{1, 2, 3})
//	// placeholders = "?,?,?"
//	// args = []any{1, 2, 3}
//
//	db.ExecContext(ctx, "delete from users where id in ("+placeholders+")", args...)
//
// SELECT 的 IN 条件请直接用 Builder 的切片展开能力：
//
//	q.Where("id in (?)", ids)
//
// 空切片返回空占位符和空参数：空的 IN 列表不是合法 SQL，调用方应跳过该语句。
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
