package tag

import (
	"reflect"
	"strings"
)

type Tag struct {
	Name  string
	Value string // 唯一值
}

func (t Tag) Hit(tag reflect.StructTag) bool {
	return tag.Get(t.Name) == t.Value
}

func (t Tag) Exist(tag reflect.StructTag) bool {
	_, ok := tag.Lookup(t.Name)
	return ok
}

func (t Tag) Contain(tag reflect.StructTag, val string) bool {
	return strings.Contains(tag.Get(t.Name), val)
}
