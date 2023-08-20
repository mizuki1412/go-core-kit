package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/spf13/cast"
)

type Role struct {
	Id          int32           `json:"id" db:"id" pk:"true" table:"role" auto:"true"`
	Department  *Department     `json:"department,omitempty" db:"department"`
	Name        class.String    `json:"name,omitempty" db:"name"`
	Description class.String    `json:"description,omitempty" db:"description"`
	Privileges  class.ArrString `json:"privileges,omitempty" db:"privileges"`
	CreateDt    class.Time      `json:"createDt,omitempty" db:"createdt"`
	Off         class.Bool      `json:"off,omitempty" db:"off" logicDel:"true"`
	Extend      class.MapString `json:"extend,omitempty" db:"extend" description:"immutable:不可删除"`
}

func (th *Role) Scan(value any) error {
	if value == nil {
		return nil
	}
	id := cast.ToInt32(value)
	th.Id = id
	return nil
}
func (th *Role) Value() (driver.Value, error) {
	return int64(th.Id), nil
}

type RoleList []*Role

func (l RoleList) Len() int           { return len(l) }
func (l RoleList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l RoleList) Less(i, j int) bool { return l[i].Id < l[j].Id }
func (l RoleList) Filter(fun func(ele *Role) bool) RoleList {
	arr := make(RoleList, 0, len(l))
	for _, e := range l {
		if fun(e) {
			arr = append(arr, e)
		}
	}
	return arr
}
func (l RoleList) Find(fun func(ele *Role) bool) *Role {
	for _, e := range l {
		if fun(e) {
			return e
		}
	}
	return nil
}
