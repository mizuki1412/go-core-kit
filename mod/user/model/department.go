package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/spf13/cast"
)

type Department struct {
	Id       int64           `auto:"true" json:"id" db:"id" pk:"true" table:"sys_department"`
	No       class.String    `json:"no,omitempty" db:"no" comment:"编号"`
	Name     class.String    `json:"name,omitempty" db:"name"`
	Descr    class.String    `json:"descr,omitempty" db:"descr" comment:"描述"`
	Parent   *Department     `json:"parent,omitempty" db:"parent"`
	Extend   class.MapString `json:"extend,omitempty" db:"extend"`
	CreateDt class.Time      `json:"createDt,omitempty" db:"createdt"`
	Deleted  class.Bool      `json:"-" db:"deleted" logicDel:"true"`
	Children []*Department   `json:"children"`
}

func (th *Department) Scan(value any) error {
	if value == nil {
		return nil
	}
	id := cast.ToInt64(value)
	th.Id = id
	return nil
}
func (th *Department) Value() (driver.Value, error) {
	return th.Id, nil
}
