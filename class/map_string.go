package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/class/utils"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/library/mapkit"
	"github.com/spf13/cast"
)

type MapString struct {
	Map   map[string]any
	Valid bool
}

func (th MapString) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.Marshal(th.Map)
	}
	return []byte("null"), nil
}

func (th *MapString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	var s map[string]any
	if err := jsonkit.Unmarshal(data, &s); err != nil {
		return err
	}
	th.Valid = true
	th.Map = s
	return nil
}

// Scan implements the Scanner interface.
func (th *MapString) Scan(value any) error {
	if value == nil {
		th.Map, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	th.Map = jsonkit.ParseMap(utils.TransScanValue2String(value))
	// todo no error
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

func NewMapString(val any) MapString {
	th := MapString{}
	if val != nil {
		th.Set(val)
	} else {
		th.Set(map[string]any{})
	}
	return th
}
func NMapString(val any) *MapString {
	th := &MapString{}
	if val != nil {
		th.Set(val)
	} else {
		th.Set(map[string]any{})
	}
	return th
}

func (th *MapString) Set(val any) {
	switch val.(type) {
	case MapString:
		v := val.(MapString)
		if v.Map == nil {
			th.Map = map[string]any{}
		} else {
			th.Map = v.Map
		}
		th.Valid = v.Valid
	case *MapString:
		v := val.(*MapString)
		if v.Map == nil {
			th.Map = map[string]any{}
		} else {
			th.Map = v.Map
		}
		th.Valid = v.Valid
	case map[string]any:
		v := val.(map[string]any)
		th.Map = v
		th.Valid = true
	default:
		panic(exception.New("class.MapString set error"))
	}
}

func (th *MapString) PutAll(val map[string]any) {
	if th.Map == nil {
		th.Map = map[string]any{}
	}
	mapkit.PutAll(th.Map, val)
	th.Valid = true
}

func (th *MapString) PutIfAbsent(key string, val any) {
	if th.Map == nil {
		th.Map = map[string]any{}
	}
	if _, ok := th.Map[key]; !ok {
		th.Map[key] = val
	}
	th.Valid = true
}

func (th *MapString) Put(key string, val any) {
	if th.Map == nil {
		th.Map = map[string]any{}
	}
	th.Map[key] = val
	th.Valid = true
}

func (th *MapString) Remove() {
	th.Valid = false
	clear(th.Map)
}

func (th MapString) IsEmpty() bool {
	if !th.Valid {
		return true
	}
	if len(th.Map) == 0 {
		return true
	}
	return false
}

func (th MapString) Contains(key string) bool {
	v, ok := th.Map[key]
	if ok {
		return v != nil
	}
	return ok
}

func (th MapString) GetOrDefault(key string, d any) any {
	v, ok := th.Map[key]
	if ok {
		return v
	}
	return d
}

func (th MapString) Get(key string) any {
	v, ok := th.Map[key]
	if ok {
		return v
	}
	return nil
}
func (th MapString) GetString(key string) string {
	return cast.ToString(th.Get(key))
}
func (th MapString) GetInt32(key string) int32 {
	return cast.ToInt32(th.Get(key))
}
func (th MapString) GetFloat64(key string) float64 {
	return cast.ToFloat64(th.Get(key))
}
func (th MapString) GetBool(key string) bool {
	return cast.ToBool(th.Get(key))
}
func (th MapString) GetMap(key string) map[string]any {
	return cast.ToStringMap(th.Get(key))
}
