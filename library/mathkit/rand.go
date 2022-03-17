package mathkit

import (
	"github.com/spf13/cast"
	"math/rand"
)

func RandFloat64(min, max any) float64 {
	v1 := cast.ToFloat64(min)
	v2 := cast.ToFloat64(max)
	if v1 > v2 {
		t := v1
		v1 = v2
		v2 = t
	}
	return rand.Float64()*(v2-v1) + v1
}

func RandInt32(min, max any) int32 {
	v1 := cast.ToInt32(min)
	v2 := cast.ToInt32(max)
	if v1 > v2 {
		t := v1
		v1 = v2
		v2 = t
	}
	return cast.ToInt32(rand.Float64()*cast.ToFloat64(v2-v1)) + v1
}
