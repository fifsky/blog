package store

import (
	"strings"
)

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
