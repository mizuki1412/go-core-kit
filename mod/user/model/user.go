package model

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/spf13/cast"
)

type User struct {
	Id       int32           `json:"id,omitempty" db:"id" pk:"true" tablename:"admin_user" autoincrement:"true"`
	Role     *Role           `json:"role,omitempty" db:"role"`
	Username class.String    `json:"username,omitempty" db:"username"`
	Name     class.String    `json:"name,omitempty" db:"name"`
	Phone    class.String    `json:"phone,omitempty" db:"phone"`
	Pwd      class.String    `json:"-" db:"pwd"`
	Gender   class.Int32     `json:"gender" db:"gender"`
	Image    class.String    `json:"image,omitempty" db:"image"`
	Address  class.String    `json:"address,omitempty" db:"address"`
	Off      class.Int32     `json:"off" db:"off"`
	Extend   class.MapString `json:"extend,omitempty" db:"extend" description:""`
	CreateDt class.Time      `json:"createDt,omitempty" db:"createdt"`
}

const UserOffOK = 0
const UserOffFreeze = 1
const UserOffDelete = -1

func (th *User) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	id := cast.ToInt32(value)
	th.Id = id
	return nil
}

func (th User) Value() (driver.Value, error) {
	return int64(th.Id), nil
}

// 判断属于某个部门
func (th *User) BelongDepartment(department int32) bool {
	return th != nil && th.Role != nil && th.Role.Department != nil && th.Role.Department.Id == department
}

// 判断是否有权限
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
func (l UserList) MapReduce(fun func(ele *User) interface{}) []interface{} {
	var results []interface{}
	for _, e := range l {
		results = append(results, fun(e))
	}
	return results
}
