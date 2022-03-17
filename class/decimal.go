package class

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type Decimal struct {
	decimal.NullDecimal
}

func NewDecimal(val any) *Decimal {
	th := &Decimal{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func ConstDecimal(v int32) decimal.Decimal {
	return decimal.NewFromInt32(v)
}

func (th *Decimal) Set(val any) *Decimal {
	if v, ok := val.(decimal.Decimal); ok {
		th.Valid = true
		val = v.String()
	} else if v, ok := val.(Decimal); ok {
		th.Valid = v.Valid
		val = v.Decimal.String()
	} else if v, ok := val.(*Decimal); ok {
		th.Valid = v.Valid
		val = v.Decimal.String()
	}
	v, err := decimal.NewFromString(cast.ToString(val))
	if err == nil {
		th.Valid = true
		th.Decimal = v
	} else {
		panic(exception.New("class.Decimal set error: " + err.Error()))
	}
	return th
}

func (th *Decimal) Round(place int32) *Decimal {
	th.Decimal = th.Decimal.Round(place)
	return th
}

func (th *Decimal) DivRound(d2 *Decimal, place int32) *Decimal {
	// 对0的处理
	if d2 == nil || d2.Float64() == 0 {
		th.Decimal = decimal.NewFromInt32(0)
	} else {
		th.Decimal = th.Decimal.DivRound(d2.Decimal, place)
	}
	return th
}

func (th *Decimal) Mul(d2 *Decimal) *Decimal {
	temp := NewDecimal(th.Decimal)
	temp.Decimal = temp.Decimal.Mul(d2.Decimal)
	return temp
}

func (th *Decimal) Add(d2 *Decimal) *Decimal {
	temp := NewDecimal(th.Decimal)
	temp.Decimal = temp.Decimal.Add(d2.Decimal)
	return temp
}

func (th *Decimal) Sub(d2 *Decimal) *Decimal {
	temp := NewDecimal(th.Decimal)
	temp.Decimal = temp.Decimal.Sub(d2.Decimal)
	return temp
}

func (th *Decimal) Div(d2 *Decimal) *Decimal {
	temp := NewDecimal(th.Decimal)
	temp.Decimal = temp.Decimal.Div(d2.Decimal)
	return temp
}

func (th Decimal) IsValid() bool {
	return th.Valid
}

func (th *Decimal) Float64() float64 {
	val, _ := th.Decimal.Float64()
	return val
}
