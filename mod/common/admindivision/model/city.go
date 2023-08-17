package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class"
)

type City struct {
	Code     class.String `json:"code,omitempty" db:"code" pk:"true" table:"city"`
	Name     class.String `json:"name,omitempty" db:"name"`
	Province *Province    `json:"province,omitempty" db:"province"`
}

func (th *City) Scan(value any) error {
	if value == nil {
		return nil
	}
	th.Code.Set(value)
	return nil
}
func (th City) Value() (driver.Value, error) {
	return th.Code.String, nil
}
