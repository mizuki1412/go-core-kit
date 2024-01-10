package class

import (
	"database/sql/driver"
	"errors"
	"github.com/lib/pq"
	"github.com/mizuki1412/go-core-kit/v2/class/const/sqlconst"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/spf13/cast"
)

// ArrInt 针对PG的array，或 mysql 里用 json
type ArrInt struct {
	Array     []int64
	Valid     bool
	dbDriver  string // 指定数据库型号
	forceJson bool   // 考虑到 pg 时可能 array or jsonb，true 则强制 jsonb
}

func (th *ArrInt) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.Marshal(th.Array)
	}
	return []byte("null"), nil
}
func (th *ArrInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	var s []int64
	if err := jsonkit.Unmarshal(data, &s); err != nil {
		return err
	}
	th.Valid = true
	th.Array = s
	return nil
}

func (th *ArrInt) SetDBDriver(driver string) {
	th.dbDriver = driver
}
func (th *ArrInt) SetForceJson(f bool) {
	th.forceJson = f
}

// Scan implements the Scanner interface.
func (th *ArrInt) Scan(value any) error {
	if value == nil || len(value.([]byte)) == 0 {
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
		var v pq.Int64Array = th.Array
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
func (th ArrInt) Value() (driver.Value, error) {
	if !th.Valid {
		return nil, nil
	}
	if th.dbDriver == sqlconst.Postgres && !th.forceJson {
		var v pq.Int64Array = th.Array
		return v.Value()
	} else {
		return jsonkit.ToString(th.Array), nil
	}
}

func NewArrInt(val any) ArrInt {
	th := ArrInt{}
	if val != nil {
		th.Set(val)
	}
	return th
}
func NArrInt(val any) *ArrInt {
	th := &ArrInt{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th ArrInt) IsValid() bool {
	return th.Valid
}

func (th *ArrInt) Set(val any) *ArrInt {
	if v, ok := val.([]int64); ok {
		th.Array = v
	} else if v, ok := val.(ArrInt); ok {
		th.Array = v.Array
	} else if v, ok := val.([]int64); ok {
		th.Array = v
	} else if v, ok := val.([]int32); ok {
		arr := make([]int64, 0, len(v))
		for _, e := range v {
			arr = append(arr, cast.ToInt64(e))
		}
		th.Array = arr
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

func (th *ArrInt) Add(vals ...any) *ArrInt {
	for _, e := range vals {
		th.Array = append(th.Array, cast.ToInt64(e))
	}
	th.Valid = true
	return th
}

func (th *ArrInt) Contains(val any) bool {
	for _, e := range th.Array {
		if e == cast.ToInt64(val) {
			return true
		}
	}
	return false
}
