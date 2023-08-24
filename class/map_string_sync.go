package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/mapkit"
	"sync"
)

// MapStringSync 同时继承scan和value方法
type MapStringSync struct {
	sync.RWMutex
	Map   map[string]any
	Valid bool
}

// todo 序列化时暂无加锁

func (th *MapStringSync) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.Marshal(th.Map)
	}
	return []byte("null"), nil
}

func (th *MapStringSync) UnmarshalJSON(data []byte) error {
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
func (th *MapStringSync) Scan(value any) error {
	if value == nil {
		th.Map, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	th.Map = jsonkit.ParseMap(string(value.([]byte)))
	return nil
}

// Value implements the driver Valuer interface.
func (th *MapStringSync) Value() (driver.Value, error) {
	if !th.Valid || th.Map == nil {
		return nil, nil
	}
	return jsonkit.ToString(th.Map), nil
}

func (th *MapStringSync) IsValid() bool {
	return th.Valid
}

func NewMapStringSync(val any) MapStringSync {
	th := MapStringSync{}
	if val != nil {
		th.Set(val)
	} else {
		th.Set(map[string]any{})
	}
	return th
}
func NMapStringSync(val any) *MapStringSync {
	th := &MapStringSync{}
	if val != nil {
		th.Set(val)
	} else {
		th.Set(map[string]any{})
	}
	return th
}

func (th *MapStringSync) Set(val any) {
	th.Lock()
	defer th.Unlock()
	switch val.(type) {
	case MapStringSync:
		if val.(MapStringSync).Map == nil {
			th.Map = map[string]any{}
		} else {
			th.Map = val.(MapStringSync).Map
		}
		th.Valid = val.(MapStringSync).Valid
	case *MapStringSync:
		if val.(*MapStringSync).Map == nil {
			th.Map = map[string]any{}
		} else {
			th.Map = val.(*MapStringSync).Map
		}
		th.Valid = val.(*MapStringSync).Valid
	case MapString:
		if val.(MapString).Map == nil {
			th.Map = map[string]any{}
		} else {
			th.Map = val.(MapString).Map
		}
		th.Valid = val.(MapString).Valid
	case *MapString:
		if val.(*MapString).Map == nil {
			th.Map = map[string]any{}
		} else {
			th.Map = val.(*MapString).Map
		}
		th.Valid = val.(*MapString).Valid
	case map[string]any:
		v := val.(map[string]any)
		th.Map = v
		th.Valid = true
	default:
		panic(exception.New("class.MapStringSync set error"))
	}
}

func (th *MapStringSync) PutAll(val map[string]any) {
	th.Lock()
	defer th.Unlock()
	if th.Map == nil {
		th.Map = map[string]any{}
	}
	mapkit.PutAll(th.Map, val)
	th.Valid = true
}

func (th *MapStringSync) PutIfAbsent(key string, val any) {
	th.Lock()
	defer th.Unlock()
	if th.Map == nil {
		th.Map = map[string]any{}
	}
	if _, ok := th.Map[key]; !ok {
		th.Map[key] = val
	}
	th.Valid = true
}

func (th *MapStringSync) Put(key string, val any) {
	th.Lock()
	defer th.Unlock()
	if th.Map == nil {
		th.Map = map[string]any{}
	}
	th.Map[key] = val
	th.Valid = true
}

func (th *MapStringSync) Remove() {
	th.Lock()
	defer th.Unlock()
	th.Valid = false
	th.Map = map[string]any{}
}

func (th *MapStringSync) Delete(key string) {
	th.Lock()
	defer th.Unlock()
	delete(th.Map, key)
}

func (th *MapStringSync) IsEmpty() bool {
	if !th.Valid {
		return true
	}
	if len(th.Map) == 0 {
		return true
	}
	return false
}

func (th *MapStringSync) Contains(key string) bool {
	th.RLock()
	defer th.RUnlock()
	v, ok := th.Map[key]
	if ok {
		return v != nil
	}
	return ok
}

func (th *MapStringSync) Get(key string) any {
	th.RLock()
	defer th.RUnlock()
	v, _ := th.Map[key]
	return v
}

func (th *MapStringSync) Entries() map[string]any {
	th.RLock()
	defer th.RUnlock()
	m := map[string]any{}
	for k, v := range th.Map {
		if vv, ok := v.(*MapStringSync); ok {
			m[k] = vv.Entries()
		} else {
			m[k] = v
		}
	}
	return m
}
