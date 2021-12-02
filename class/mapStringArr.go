package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
)

/** 针对PG的jsonb，其中是array形式的 */

type MapStringArr struct {
	Arr   []map[string]interface{}
	Valid bool
}

func (th MapStringArr) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Arr)
	}
	return []byte("null"), nil
}
func (th *MapStringArr) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	var s []map[string]interface{}
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	th.Valid = true
	th.Arr = s
	return nil
}

// Scan implements the Scanner interface.
func (th *MapStringArr) Scan(value interface{}) error {
	if value == nil {
		th.Arr, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	var s []map[string]interface{}
	err := jsonkit.ParseObj(string(value.([]byte)), &s)
	if err != nil {
		logkit.Error(exception.New(err.Error()))
	}
	th.Arr = s
	return nil
}

// Value implements the driver Valuer interface.
func (th MapStringArr) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return jsonkit.ToString(th.Arr), nil
}

func (th MapStringArr) IsValid() bool {
	return th.Valid
}

func NewMapStringArr(val interface{}) *MapStringArr {
	th := &MapStringArr{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *MapStringArr) Set(val interface{}) *MapStringArr {
	if v, ok := val.(MapStringArr); ok {
		th.Arr = v.Arr
		th.Valid = v.Valid
	} else if v, ok := val.([]map[string]interface{}); ok {
		th.Arr = v
		th.Valid = true
	} else {
		panic(exception.New("class.MapStringArr set error"))
	}
	return th
}

func (th *MapStringArr) Length() int {
	return len(th.Arr)
}
func (th *MapStringArr) Remove() *MapStringArr {
	th.Valid = false
	th.Arr = []map[string]interface{}{}
	return th
}

func (th *MapStringArr) IsEmpty() bool {
	if !th.Valid {
		return true
	}
	if len(th.Arr) == 0 {
		return true
	}
	return false
}
