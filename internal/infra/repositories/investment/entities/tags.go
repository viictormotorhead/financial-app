package entities

import (
	"database/sql/driver"
	"fmt"

	"github.com/lib/pq"
)

// StringList maps to a PostgreSQL text[] column.
type StringList []string

func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return pq.Array([]string{}).Value()
	}
	return pq.Array([]string(s)).Value()
}

func (s *StringList) Scan(value interface{}) error {
	if value == nil {
		*s = StringList{}
		return nil
	}

	var arr pq.StringArray

	switch v := value.(type) {
	case []string:
		*s = StringList(v)
		return nil
	case pq.StringArray:
		*s = StringList(v)
		return nil
	case string:
		if err := arr.Scan(v); err != nil {
			return err
		}
		*s = StringList(arr)
		return nil
	case []byte:
		if err := arr.Scan(string(v)); err != nil {
			return err
		}
		*s = StringList(arr)
		return nil
	default:
		if err := arr.Scan(value); err != nil {
			return fmt.Errorf("scan tags: unsupported type %T: %w", value, err)
		}
		*s = StringList(arr)
		return nil
	}
}

func (s StringList) Strings() []string {
	if s == nil {
		return []string{}
	}
	return []string(s)
}
