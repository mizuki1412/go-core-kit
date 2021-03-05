package mathkit

import (
	"github.com/spf13/cast"
	"math/rand"
)

func RandFloat64(min, max interface{}) float64 {
	v1 := cast.ToFloat64(min)
	v2 := cast.ToFloat64(max)
	if v1 > v2 {
		t := v1
		v1 = v2
		v2 = t
	}
	return rand.Float64()*(v2-v1) + v1
}
