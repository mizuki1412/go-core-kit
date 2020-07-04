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
	v, err := decimal.NewFromString(cast.ToString(val))
	if err != nil {
		th.Valid = true
		th.Decimal = v
	} else {
		panic(exception.New("decimal set error"))
	}
}
