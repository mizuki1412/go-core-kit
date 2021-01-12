package timekit

import "time"

// 是否存在交集，不含点
func IsOverlap(base []time.Time, layout []time.Time) bool {
	if len(base) != 2 || len(layout) != 2 {
		return false
	}
	return base[1].Unix() > layout[0].Unix() && base[0].Unix() < layout[1].Unix()
}

// 取base中不含layout的部分
func GetNotOverlapArray(base []time.Time, layout []time.Time) ([]time.Time, []time.Time) {
	base1 := base[0].Unix()
	base2 := base[1].Unix()
	layout1 := layout[0].Unix()
	layout2 := layout[1].Unix()
	if layout1 >= base2 || layout2 <= base1 {
		return base, nil
	}
	if layout2 < base2 && layout1 < base1 {
		return []time.Time{layout[1], base[1]}, nil
	} else if layout1 > base1 && layout2 > base2 {
		return []time.Time{base[0], layout[0]}, nil
	} else if layout1 > base1 && layout2 < base2 {
		return []time.Time{base[0], layout[0]}, []time.Time{layout[1], base[1]}
	} else {
		return nil, nil
	}
}

// 合并两组时间，返回nil则未合并
func MergeArray(base []time.Time, layout []time.Time) []time.Time {
	base1 := base[0].Unix()
	base2 := base[1].Unix()
	layout1 := layout[0].Unix()
	layout2 := layout[1].Unix()
	if layout1 >= base2 || layout2 <= base1 {
		return nil
	}
	if layout2 < base2 && layout1 < base1 {
		return []time.Time{layout[0], base[1]}
	} else if layout1 > base1 && layout2 > base2 {
		return []time.Time{base[0], layout[1]}
	} else if layout1 > base1 && layout2 < base2 {
		return base
	} else {
		return layout
	}
}

func GetNotOverlapSegment(base []time.Time, layouts [][]time.Time) [][]time.Time {
	// todo
	return nil
}
