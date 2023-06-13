package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/spf13/cast"
)

type PrivilegeConstant struct {
	Id   string       `json:"id" db:"id" pk:"true" tablename:"privilege_constant"`
	Name class.String `json:"name,omitempty" db:"name"`
	Type class.String `json:"type,omitempty" db:"type" description:"暂不用"`
	Sort int32        `json:"sort" db:"sort"`
}

func (th *PrivilegeConstant) Scan(value any) error {
	if value == nil {
		return nil
	}
	id := cast.ToString(value)
	th.Id = id
	return nil
}
func (th PrivilegeConstant) Value() (driver.Value, error) {
	return th.Id, nil
}
