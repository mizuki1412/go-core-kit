package class

import (
	"database/sql/driver"
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
	return jsonkit.JSON().Marshal(nil)
}
func (th *MapString) UnmarshalJSON(data []byte) error {
	var s *map[string]interface{}
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		th.Valid = true
		th.Map = *s
	} else {
		th.Valid = false
	}
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

func (th *MapString) Set(val map[string]interface{}) {
	th.Map = val
	th.Valid = true
}

func (th *MapString) PutAll(val map[string]interface{}) {
	mapkit.PutAll(th.Map, val)
	th.Valid = true
}
