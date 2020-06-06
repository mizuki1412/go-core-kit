package class

import (
	"database/sql/driver"
	"mizuki/framework/core-kit/library/jsonkit"
	"mizuki/framework/core-kit/library/stringkit"
	"strings"
)

/**
针对PG的array
*/

// 同时继承scan和value方法
type ArrString struct {
	Array []string
	Valid bool
}

func (th ArrString) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Array)
	}
	return jsonkit.JSON().Marshal(nil)
}
func (th *ArrString) UnmarshalJSON(data []byte) error {
	var s *[]string
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
func (th *ArrString) Scan(value interface{}) error {
	if value == nil {
		th.Array, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	// parse eg: {abc,qq} to ["abc","qq"]
	bytes := value.([]byte)
	if len(bytes) <= 2 {
		th.Array = []string{}
		return nil
	}
	th.Array = strings.Split(string(bytes[1:len(bytes)-1]), ",")
	return nil
}

// Value implements the driver Valuer interface.
func (th ArrString) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	return "{" + stringkit.ConcatWith(th.Array, ",", "'") + "}", nil
}

func (th *ArrString) Set(val []string) {
	th.Array = val
	th.Valid = true
}
