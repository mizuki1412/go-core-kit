package class

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type Decimal struct {
	decimal.NullDecimal
}

func NewDecimal(val interface{}) *Decimal {
	th := &Decimal{}
	if val != nil {
		th.Set(val)
	}
	return th
}

func (th *Decimal) Set(val interface{}) *Decimal {
	if v, ok := val.(decimal.Decimal); ok {
		th.Valid = true
		th.Decimal = v
	} else if v, ok := val.(Decimal); ok {
		th.Valid = true
		th.Decimal = v.Decimal
	} else {
		v, err := decimal.NewFromString(cast.ToString(val))
		if err == nil {
			th.Valid = true
			th.Decimal = v
		} else {
			panic(exception.New("class.Decimal set error: " + err.Error()))
		}
	}
	return th
}

func (th *Decimal) Round(place int32) *Decimal {
	th.Decimal = th.Decimal.Round(place)
	return th
}

func (th Decimal) IsValid() bool {
	return th.Valid
}

func (th *Decimal) Float64() float64 {
	val, _ := th.Decimal.Float64()
	return val
}
