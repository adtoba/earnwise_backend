package models

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	// Convert to PostgreSQL array literal format
	result := "{"
	for i, v := range a {
		if i > 0 {
			result += ","
		}
		result += "\"" + v + "\""
	}
	result += "}"
	return result, nil
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// Remove the curly braces and split by comma
		str := string(v)
		str = strings.Trim(str, "{}")
		if str == "" {
			*a = StringArray{}
			return nil
		}
		// Split and remove quotes
		parts := strings.Split(str, ",")
		result := make(StringArray, len(parts))
		for i, part := range parts {
			result[i] = strings.Trim(part, "\"")
		}
		*a = result
		return nil
	case string:
		return a.Scan([]byte(v))
	default:
		return errors.New("type assertion to []byte failed")
	}
}
