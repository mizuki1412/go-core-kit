package class

import (
	"database/sql/driver"
	"github.com/lib/pq"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
)

type ArrString struct {
	Array pq.StringArray
	Valid bool
}

func (th ArrString) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.Marshal(th.Array)
	}
	return []byte("null"), nil
}
func (th *ArrString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	var s pq.StringArray
	if err := jsonkit.Unmarshal(data, &s); err != nil {
		return err
	}
	th.Valid = true
	th.Array = s
	return nil
}

func (th ArrString) IsValid() bool {
	return th.Valid
}

// Scan implements the Scanner interface.
func (th *ArrString) Scan(value any) error {
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

func NewArrString(val any) ArrString {
	th := ArrString{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *ArrString) Set(val any) {
	if v, ok := val.([]string); ok {
		th.Array = v
	} else if v, ok := val.(ArrString); ok {
		th.Array = v.Array
	} else if v, ok := val.(pq.StringArray); ok {
		th.Array = v
	} else {
		panic(exception.New("class.ArrString set error"))
	}
	th.Valid = true
}

func (th *ArrString) Length() int {
	return len(th.Array)
}
func (th *ArrString) Remove() {
	th.Valid = false
	th.Array = []string{}
}

func (th *ArrString) Delete(vals ...string) {
	temp := th.Array
	for _, e := range vals {
		temp = arraykit.StringDelete(temp, e)
	}
	th.Array = temp
}

func (th *ArrString) Add(vals ...string) {
	th.Array = append(th.Array, vals...)
	th.Valid = true
}

func (th *ArrString) Contains(val string) bool {
	for _, e := range th.Array {
		if e == val {
			return true
		}
	}
	return false
}
