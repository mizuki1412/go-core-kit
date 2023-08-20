package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class"
)

type Province struct {
	Code   class.String `json:"code,omitempty" db:"code" pk:"true" table:"province"`
	Name   class.String `json:"name,omitempty" db:"name"`
	Cities []*City      `json:"cities"`
}

func (th *Province) Scan(value any) error {
	if value == nil {
		return nil
	}
	th.Code.Set(value)
	return nil
}
func (th *Province) Value() (driver.Value, error) {
	return th.Code.String, nil
}
