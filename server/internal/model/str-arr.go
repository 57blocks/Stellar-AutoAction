package model

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type StrList []string

func (a *StrList) Value() (driver.Value, error) {
	return "{" + strings.Join(*a, ",") + "}", nil
}

func (a *StrList) Scan(value interface{}) error {
	if value == nil {
		*a = StrList{}
		return nil
	}

	s, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}

	*a = strings.Split(strings.Trim(s, "{}"), ",")

	return nil
}
