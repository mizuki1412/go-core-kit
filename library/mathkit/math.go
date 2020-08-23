package mathkit

import (
	"github.com/shopspring/decimal"
	"math"
)

func GroupNum(sum int, group int) int {
	if sum == 0 {
		return 1
	}
	if sum%group == 0 {
		return sum / group
	} else {
		return sum/group + 1
	}
}

// 保留小数
func FloatRound(val float64, num int32) float64 {
	if math.IsNaN(val) {
		return 0
	}
	f, _ := decimal.NewFromFloat(val).Round(num).Float64()
	return f
}
