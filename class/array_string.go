package class

import (
	"database/sql/driver"
	"errors"
	"github.com/lib/pq"
	"github.com/mizuki1412/go-core-kit/v2/class/const/sqlconst"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/arraykit"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
)

type ArrString struct {
	Array     []string
	Valid     bool
	dbDriver  string // 指定数据库型号
	forceJson bool   // 考虑到 pg 时可能 array or jsonb，true 则强制 jsonb
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

func (th *ArrString) SetDBDriver(driver string) {
	th.dbDriver = driver
}
func (th *ArrString) SetForceJson(f bool) {
	th.forceJson = f
}

// Scan implements the Scanner interface.
func (th *ArrString) Scan(value any) error {
	if value == nil {
		th.Array, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	// 通过首字符判断
	val := string(value.([]byte))
	switch val[0] {
	case '[':
		return jsonkit.ParseObj(val, &th.Array)
	case '{':
		var v pq.StringArray = th.Array
		err := v.Scan(value)
		if err != nil {
			return err
		}
		th.Array = v
		return nil
	default:
		return errors.New("scan not support")
	}
}

// Value implements the driver Valuer interface.
func (th ArrString) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	if th.dbDriver == sqlconst.Postgres && !th.forceJson {
		var v pq.StringArray = th.Array
		return v.Value()
	} else {
		return jsonkit.ToString(th.Array), nil
	}
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
