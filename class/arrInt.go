package class

import (
	"database/sql/driver"
	"github.com/lib/pq"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
)

/**
针对PG的array
*/

// 同时继承scan和value方法
type ArrInt struct {
	Array pq.Int64Array
	Valid bool
}

func (th ArrInt) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Array)
	}
	return []byte("null"), nil
}
func (th *ArrInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	var s []int64
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	th.Valid = true
	th.Array = s
	return nil
}

// Scan implements the Scanner interface.
func (th *ArrInt) Scan(value interface{}) error {
	if value == nil {
		th.Array, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	return th.Array.Scan(value)
}

// Value implements the driver Valuer interface.
func (th ArrInt) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return th.Array.Value()
}

func (th *ArrInt) Set(val []int64) {
	th.Array = val
	th.Valid = true
}

func (th *ArrInt) Length() int {
	return len(th.Array)
}
