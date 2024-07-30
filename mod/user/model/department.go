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

type DeptList []*Department

func (l DeptList) Len() int           { return len(l) }
func (l DeptList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l DeptList) Less(i, j int) bool { return l[i].Id < l[j].Id }
func (l DeptList) Filter(fun func(ele *Department) bool) DeptList {
	arr := make(DeptList, 0, len(l))
	for _, e := range l {
		if fun(e) {
			arr = append(arr, e)
		}
	}
	return arr
}
func (l DeptList) Find(fun func(ele *Department) bool) *Department {
	for _, e := range l {
		if fun(e) {
			return e
		}
	}
	return nil
}
func (l DeptList) MapReduce(fun func(ele *Department) any) []any {
	var results []any
	for _, e := range l {
		results = append(results, fun(e))
	}
	return results
}
