package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/spf13/cast"
)

// 同时继承scan和value方法
type String struct {
	String string
	Valid  bool
}

func (th String) MarshalJSON() ([]byte, error) {
	if th.Valid {
		// 可能存在逃逸字符
		return jsonkit.JSON().Marshal(th.String)
	}
	// 返回json中的null
	return []byte("null"), nil
	//return nil,nil
}
func (th *String) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	th.String = utils.UnquoteIfQuoted(data)
	th.Valid = true
	return nil
}
func (th *String) Scan(value interface{}) error {
	if value == nil {
		th.String, th.Valid = "", false
		return nil
	}
	th.Valid = true
	th.String = cast.ToString(value)
	return nil
}

// Value implements the driver Valuer interface.
func (th String) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return th.String, nil
}

func (th String) IsValid() bool {
	return th.Valid
}

func NewString(val interface{}) *String {
	th := &String{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *String) Set(val interface{}) *String {
	if v, ok := val.(String); ok {
		th.String = v.String
		th.Valid = v.Valid
	} else {
		s, err := cast.ToStringE(val)
		if err == nil {
			th.String = s
			th.Valid = true
		}
	}
	return th
}

func (th *String) Remove() *String {
	th.Valid = false
	th.String = ""
	return th
}
