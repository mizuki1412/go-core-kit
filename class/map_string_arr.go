package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/class/utils"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
)

/** 针对PG的jsonb，其中是array形式的 */

type MapStringArr struct {
	Arr   []map[string]any
	Valid bool
}

func (th MapStringArr) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.Marshal(th.Arr)
	}
	return []byte("null"), nil
}
func (th *MapStringArr) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	var s []map[string]any
	if err := jsonkit.Unmarshal(data, &s); err != nil {
		return err
	}
	th.Valid = true
	th.Arr = s
	return nil
}

// Scan implements the Scanner interface.
func (th *MapStringArr) Scan(value any) error {
	if value == nil {
		th.Arr, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	var s []map[string]any
	val := utils.TransScanValue2String(value)
	err := jsonkit.ParseObj(val, &s)
	th.Arr = s
	return err
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

func NewMapStringArr(val any) MapStringArr {
	th := MapStringArr{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *MapStringArr) Set(val any) {
	switch val.(type) {
	case MapStringArr:
		v := val.(MapStringArr)
		th.Arr = v.Arr
		th.Valid = v.Valid
	case *MapStringArr:
		v := val.(*MapStringArr)
		th.Arr = v.Arr
		th.Valid = v.Valid
	case []map[string]any:
		v := val.([]map[string]any)
		th.Arr = v
		th.Valid = true
	default:
		panic(exception.New("class.MapStringArr set error"))
	}
}

func (th MapStringArr) Length() int {
	return len(th.Arr)
}
func (th MapStringArr) Remove() {
	th.Valid = false
	th.Arr = []map[string]any{}
}

func (th MapStringArr) IsEmpty() bool {
	if !th.Valid {
		return true
	}
	if len(th.Arr) == 0 {
		return true
	}
	return false
}
