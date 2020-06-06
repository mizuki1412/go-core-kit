package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/spf13/cast"
	"strings"
)

/**
针对PG的array
*/

// 同时继承scan和value方法
type ArrInt struct {
	Array []int32
	Valid bool
}

func (th ArrInt) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Array)
	}
	return jsonkit.JSON().Marshal(nil)
}
func (th *ArrInt) UnmarshalJSON(data []byte) error {
	var s *[]int32
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		th.Valid = true
		th.Array = *s
	} else {
		th.Valid = false
	}
	return nil
}

// Scan implements the Scanner interface.
func (th *ArrInt) Scan(value interface{}) error {
	if value == nil {
		th.Array, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	// parse eg: {1,2} to [1,2]
	bytes := value.([]byte)
	if len(bytes) <= 2 {
		th.Array = []int32{}
		return nil
	}
	es := strings.Split(string(bytes[1:len(bytes)-1]), ",")
	var arr []int32
	for _, v := range es {
		arr = append(arr, cast.ToInt32(v))
	}
	th.Array = arr
	return nil
}

// Value implements the driver Valuer interface.
func (th ArrInt) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return "{" + stringkit.ConcatIntWith(th.Array, ", ") + "}", nil
}

func (th *ArrInt) Set(val []int32) {
	th.Array = val
	th.Valid = true
}
