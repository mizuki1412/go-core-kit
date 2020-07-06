package class

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type Decimal struct {
	decimal.NullDecimal
}

func (th *Decimal) Set(val interface{}) {
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
			panic(exception.New("class.Decimal set error"))
		}
	}
}
