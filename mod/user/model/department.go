package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/spf13/cast"
)

type Department struct {
	Id       int32           `json:"id" db:"id" pk:"true" tablename:"department" autoincrement:"true"`
	No       class.String    `json:"no,omitempty" db:"no"`
	Name     class.String    `json:"name,omitempty" db:"name"`
	Descr    class.String    `json:"descr,omitempty" db:"descr"`
	Parent   *Department     `json:"parent,omitempty" db:"parent"`
	Extend   class.MapString `json:"extend,omitempty" db:"extend"`
	CreateDt class.Time      `json:"createDt,omitempty" db:"createdt"`
	Children []*Department   `json:"children"`
}

func (th *Department) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	id := cast.ToInt32(value)
	th.Id = id
	return nil
}
func (th Department) Value() (driver.Value, error) {
	return int64(th.Id), nil
}
