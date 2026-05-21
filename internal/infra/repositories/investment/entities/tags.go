package entities

import (
	"database/sql/driver"
	"errors"

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
	case []byte:
		return arr.Scan(string(v))
	case string:
		return arr.Scan(v)
	default:
		if err := arr.Scan(value); err == nil {
			*s = StringList(arr)
			return nil
		}
		return errors.New("unsupported tags column type")
	}
}

func (s StringList) Strings() []string {
	if s == nil {
		return []string{}
	}
	return []string(s)
}
