package class

import (
	"database/sql/driver"
	"encoding/xml"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/class/utils"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/spf13/cast"
)

// String 同时继承scan和value方法
type String struct {
	String string
	Valid  bool
}

func (th String) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// todo ?
	if th.Valid {
		return e.EncodeElement(th.String, start)
	}
	return e.EncodeElement(nil, start)
}
func (th *String) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	// Read tag content into value
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}
	if value == "" {
		th.Valid = false
	} else {
		th.Valid = true
		th.String = utils.UnquoteIfQuoted([]byte(value))
	}
	return nil
}
func (th *String) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Value == "" {
		th.Valid = false
	} else {
		th.Valid = true
		th.String = utils.UnquoteIfQuoted([]byte(attr.Value))
	}
	return nil
}

func (th String) MarshalJSON() ([]byte, error) {
	if th.Valid {
		// 可能存在逃逸字符
		return jsonkit.Marshal(th.String)
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
func (th *String) Scan(value any) error {
	if value == nil {
		th.String, th.Valid = "", false
		return nil
	}
	var err error
	th.Valid = true
	th.String, err = cast.ToStringE(value)
	return err
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

func NewString(val any) *String {
	th := &String{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *String) Set(val any) *String {
	if v, ok := val.(String); ok {
		th.String = v.String
		th.Valid = v.Valid
	} else {
		s, err := cast.ToStringE(val)
		if err == nil {
			th.String = s
			th.Valid = true
		} else {
			panic(exception.New(err.Error()))
		}
	}
	return th
}

func (th *String) Remove() *String {
	th.Valid = false
	th.String = ""
	return th
}
