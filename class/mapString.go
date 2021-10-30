package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/mapkit"
	"github.com/spf13/cast"
)

/** 针对PG的jsonb */

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
	if !th.Valid || th.Map == nil {
		return nil, nil
	}
	return jsonkit.ToString(th.Map), nil
}

func (th MapString) IsValid() bool {
	return th.Valid
}

func NewMapString(val interface{}) *MapString {
	th := &MapString{}
	if val != nil {
		th.Set(val)
	} else {
		th.Set(map[string]interface{}{})
	}
	return th
}

func (th *MapString) Set(val interface{}) *MapString {
	if v, ok := val.(MapString); ok {
		if v.Map == nil {
			th.Map = map[string]interface{}{}
		} else {
			th.Map = v.Map
		}
		th.Valid = v.Valid
	} else if v, ok := val.(map[string]interface{}); ok {
		th.Map = v
		th.Valid = true
	} else {
		panic(exception.New("class.MapString set error"))
	}
	return th
}

func (th *MapString) PutAll(val map[string]interface{}) *MapString {
	if th.Map == nil {
		th.Map = map[string]interface{}{}
	}
	mapkit.PutAll(th.Map, val)
	th.Valid = true
	return th
}

func (th *MapString) PutIfAbsent(key string, val interface{}) *MapString {
	if th.Map == nil {
		th.Map = map[string]interface{}{}
	}
	if _, ok := th.Map[key]; !ok {
		th.Map[key] = val
	}
	th.Valid = true
	return th
}

func (th *MapString) Put(key string, val interface{}) *MapString {
	if th.Map == nil {
		th.Map = map[string]interface{}{}
	}
	th.Map[key] = val
	th.Valid = true
	return th
}

func (th *MapString) Remove() *MapString {
	th.Valid = false
	th.Map = map[string]interface{}{}
	return th
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

func (th *MapString) Contains(key string) bool {
	v, ok := th.Map[key]
	if ok {
		return v != nil
	}
	return ok
}

func (th *MapString) GetOrDefault(key string, d interface{}) interface{} {
	v, ok := th.Map[key]
	if ok {
		return v
	}
	return d
}

func (th *MapString) Get(key string) interface{} {
	v, ok := th.Map[key]
	if ok {
		return v
	}
	return nil
}
func (th *MapString) GetString(key string) string {
	return cast.ToString(th.Get(key))
}
func (th *MapString) GetInt32(key string) int32 {
	return cast.ToInt32(th.Get(key))
}
func (th *MapString) GetFloat64(key string) float64 {
	return cast.ToFloat64(th.Get(key))
}
func (th *MapString) GetBool(key string) bool {
	return cast.ToBool(th.Get(key))
}
func (th *MapString) GetMap(key string) map[string]interface{} {
	return cast.ToStringMap(th.Get(key))
}
