package class

import (
	"github.com/spf13/cast"
	"sort"
	"time"
)

//
// 两个时间一组的数组
//
type TimeArrGroup [][]time.Time

type TimePoint struct {
	Time time.Time
	Id   string
	Type int32 // 1,2; 11,12-剔除的开始结束点
}
type TimePointList []*TimePoint

func (l TimePointList) Len() int           { return len(l) }
func (l TimePointList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l TimePointList) Less(i, j int) bool { return l[i].Time.Unix() < l[j].Time.Unix() }

// 累计时间长，s
func (th TimeArrGroup) Sum() int64 {
	var all int64 = 0
	for _, t := range th {
		if len(t) != 2 {
			continue
		}
		dif := t[1].Unix() - t[0].Unix()
		if dif <= 0 {
			continue
		}
		all += dif
	}
	return all
}

func _group2TimeArr(a, b TimeArrGroup) TimePointList {
	list := make(TimePointList, 0, len(a)*2+len(b)*2)
	for i, e := range a {
		if len(e) != 2 {
			continue
		}
		list = append(list, &TimePoint{
			Time: e[0],
			Id:   "a" + cast.ToString(i),
			Type: 1,
		}, &TimePoint{
			Time: e[1],
			Id:   "a" + cast.ToString(i),
			Type: 2,
		})
	}
	for i, e := range b {
		if len(e) != 2 {
			continue
		}
		list = append(list, &TimePoint{
			Time: e[0],
			Id:   "a" + cast.ToString(i),
			Type: 1,
		}, &TimePoint{
			Time: e[1],
			Id:   "a" + cast.ToString(i),
			Type: 2,
		})
	}
	sort.Sort(list)
	return list
}

// 合并，当前的时间数组和参数的时间数组合并。
func (th TimeArrGroup) Merge(obj TimeArrGroup) TimeArrGroup {
	//list := _group2TimeArr(th,obj)

	return nil
}

// 剔除，当前的时间数组剔除参数的时间数组范围
func (th TimeArrGroup) Eliminate(obj TimeArrGroup) TimeArrGroup {
	return nil
}
