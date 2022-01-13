package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class"
)

type Area struct {
	Code class.String `json:"code,omitempty" db:"code" pk:"true" tablename:"area"`
	Name class.String `json:"name,omitempty" db:"name"`
}

func (th *Area) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	th.Code.Set(value)
	return nil
}
func (th Area) Value() (driver.Value, error) {
	return th.Code.String, nil
}
