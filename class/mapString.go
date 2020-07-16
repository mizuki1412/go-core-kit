package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/mapkit"
)

/**
针对PG的jsonb
*/

// 同时继承scan和value方法
type MapString struct {
	Map   map[string]interface{}
	Valid bool
}

func (th MapString) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Map)
	}
	return []byte("null"), nil
}
func (th *MapString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	var s map[string]interface{}
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	th.Valid = true
	th.Map = s
	return nil
}

// Scan implements the Scanner interface.
func (th *MapString) Scan(value interface{}) error {
	if value == nil {
		th.Map, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	th.Map = jsonkit.ParseMap(string(value.([]byte)))
	return nil
}

// Value implements the driver Valuer interface.
func (th MapString) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return jsonkit.ToString(th.Map), nil
}

func (th *MapString) Set(val interface{}) {
	if v, ok := val.(MapString); ok {
		th.Map = v.Map
		th.Valid = true
	} else if v, ok := val.(map[string]interface{}); ok {
		th.Map = v
		th.Valid = true
	} else {
		panic(exception.New("class.MapString set error"))
	}
}

func (th *MapString) PutAll(val map[string]interface{}) {
	if th.Map == nil {
		th.Map = map[string]interface{}{}
	}
	mapkit.PutAll(th.Map, val)
	th.Valid = true
}

func (th *MapString) PutIfAbsent(key string, val interface{}) {
	if th.Map == nil {
		th.Map = map[string]interface{}{}
	}
	if _, ok := th.Map[key]; !ok {
		th.Map[key] = val
	}
	th.Valid = true
}

func (th *MapString) Put(key string, val interface{}) {
	if th.Map == nil {
		th.Map = map[string]interface{}{}
	}
	th.Map[key] = val
	th.Valid = true
}

func (th *MapString) Remove() {
	th.Valid = false
	th.Map = map[string]interface{}{}
}

func (th *MapString) IsEmpty() bool {
	if !th.Valid {
		return true
	}
	if len(th.Map) == 0 {
		return true
	}
	return false
}
