package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Tags []string

func (t *Tags) Scan(value any) error {
	if value == nil {
		*t = nil
		return nil
	}

	var raw []byte
	switch v := value.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return fmt.Errorf("scan tags: unsupported type %T", value)
	}

	if len(raw) == 0 {
		*t = nil
		return nil
	}

	var tmp []string
	if err := json.Unmarshal(raw, &tmp); err != nil {
		return err
	}
	*t = Tags(tmp)
	return nil
}

func (t Tags) Value() (driver.Value, error) {
	if t == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(t))
}
