package jsonutil

import "encoding/json"

func Encode(v any) string {
	str, _ := json.Marshal(v)
	return string(str)
}
