package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/spf13/cast"
)

type Role struct {
	Id          int64           `json:"id" db:"id" pk:"true" table:"sys_role" auto:"true"`
	Department  *Department     `json:"department,omitempty" db:"department"`
	Name        class.String    `json:"name,omitempty" db:"name"`
	Description class.String    `json:"description,omitempty" db:"description"`
	Privileges  class.ArrString `json:"privileges,omitempty" db:"privileges"`
	CreateDt    class.Time      `json:"createDt,omitempty" db:"createdt"`
	Deleted     class.Bool      `json:"-" db:"deleted" logicDel:"true"`
	Extend      class.MapString `json:"extend,omitempty" db:"extend" comment:"immutable:不可删除"`
}

func (th *Role) Scan(value any) error {
	if value == nil {
		return nil
	}
	id := cast.ToInt64(value)
	th.Id = id
	return nil
}
func (th *Role) Value() (driver.Value, error) {
	return th.Id, nil
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
