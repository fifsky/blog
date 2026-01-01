package jsonutil

import "encoding/json"

func Encode(v interface{}) string {
	str, _ := json.Marshal(v)
	return string(str)
}
