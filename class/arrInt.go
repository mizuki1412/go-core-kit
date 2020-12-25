package class

import (
	"database/sql/driver"
	"github.com/lib/pq"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/spf13/cast"
)

/** 针对PG的array */

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

func NewArrInt(val interface{}) *ArrInt {
	th := &ArrInt{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th ArrInt) IsValid() bool {
	return th.Valid
}

func (th *ArrInt) Set(val interface{}) *ArrInt {
	if v, ok := val.([]int64); ok {
		th.Array = v
	} else if v, ok := val.(ArrInt); ok {
		th.Array = v.Array
	} else if v, ok := val.(pq.Int64Array); ok {
		th.Array = v
	} else {
		panic(exception.New("class.ArrInt set error"))
	}
	th.Valid = true
	return th
}

func (th *ArrInt) Length() int {
	return len(th.Array)
}

func (th *ArrInt) Remove() *ArrInt {
	th.Valid = false
	th.Array = []int64{}
	return th
}

func (th *ArrInt) ToInt32Slice() []int32 {
	list := make([]int32, 0, len(th.Array))
	for _, e := range th.Array {
		list = append(list, cast.ToInt32(e))
	}
	return list
}

func (th *ArrInt) Add(vals ...int64) *ArrInt {
	th.Array = append(th.Array, vals...)
	th.Valid = true
	return th
}

func (th *ArrInt) Add32(vals ...int32) *ArrInt {
	for _, e := range vals {
		th.Array = append(th.Array, cast.ToInt64(e))
	}
	th.Valid = true
	return th
}
