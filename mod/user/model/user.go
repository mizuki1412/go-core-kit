package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/spf13/cast"
)

type User struct {
	Id         int64           `json:"id,omitempty" db:"id" pk:"true" table:"sys_user" auto:"true"`
	Role       *Role           `json:"role,omitempty" db:"role"`
	Department *Department     `json:"department,omitempty" db:"department"`
	Username   class.String    `json:"username,omitempty" db:"username"`
	Name       class.String    `json:"name,omitempty" db:"name"`
	Phone      class.String    `json:"phone,omitempty" db:"phone"`
	Pwd        class.String    `json:"-" db:"pwd"`
	Gender     class.Int32     `json:"gender,omitempty" db:"gender" comment:"1-nan,2-nv"`
	Image      class.String    `json:"image,omitempty" db:"image" comment:"头像"`
	Address    class.String    `json:"address,omitempty" db:"address"`
	Status     class.Int32     `json:"status,omitempty" db:"status" comment:"冻结 1"`
	Deleted    class.Bool      `json:"-" db:"deleted" logicDel:"true"`
	Extend     class.MapString `json:"extend,omitempty" db:"extend" comment:"权限剔除privilegeExclude:[], 不可删除immutable:bool"`
	CreateDt   class.Time      `json:"createDt,omitempty" db:"createdt"`
}

const UserStatusOK = 0
const UserStatusFreeze = 1

func (th *User) Scan(value any) error {
	if value == nil {
		return nil
	}
	id := cast.ToInt64(value)
	th.Id = id
	return nil
}

func (th *User) Value() (driver.Value, error) {
	return th.Id, nil
}

// BelongDepartment 判断属于某个部门
func (th *User) BelongDepartment(department int64) bool {
	return th != nil && th.Role != nil && th.Role.Department != nil && th.Role.Department.Id == department
}

// HasPrivilege 判断是否有权限
func (th *User) HasPrivilege(privilege string) bool {
	return th != nil && th.Role != nil && th.Role.Privileges.Contains(privilege)
}

type UserList []*User

func (l UserList) Len() int           { return len(l) }
func (l UserList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l UserList) Less(i, j int) bool { return l[i].Id < l[j].Id }
func (l UserList) Filter(fun func(ele *User) bool) UserList {
	arr := make(UserList, 0, len(l))
	for _, e := range l {
		if fun(e) {
			arr = append(arr, e)
		}
	}
	return arr
}
func (l UserList) Find(fun func(ele *User) bool) *User {
	for _, e := range l {
		if fun(e) {
			return e
		}
	}
	return nil
}
func (l UserList) MapReduce(fun func(ele *User) any) []any {
	var results []any
	for _, e := range l {
		results = append(results, fun(e))
	}
	return results
}
