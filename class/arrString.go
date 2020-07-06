package class

import (
	"database/sql/driver"
	"github.com/lib/pq"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
)

/**
针对PG的array
*/

// 同时继承scan和value方法
type ArrString struct {
	Array pq.StringArray
	Valid bool
}

func (th ArrString) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Array)
	}
	return []byte("null"), nil
}
func (th *ArrString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	var s pq.StringArray
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	th.Valid = true
	th.Array = s
	return nil
}

// Scan implements the Scanner interface.
func (th *ArrString) Scan(value interface{}) error {
	if value == nil {
		th.Array, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	return th.Array.Scan(value)
}

// Value implements the driver Valuer interface.
func (th ArrString) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return th.Array.Value()
}

func (th *ArrString) Set(val interface{}) {
	if v, ok := val.([]string); ok {
		th.Array = v
		th.Valid = true
	} else if v, ok := val.(ArrString); ok {
		th.Array = v.Array
		th.Valid = true
	} else {
		panic(exception.New("class.ArrString set error"))
	}
}

func (th *ArrString) Length() int {
	return len(th.Array)
}
