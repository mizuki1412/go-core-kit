package class

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/spf13/cast"
	"sort"
)

type Int64List []int64

func (l Int64List) Len() int           { return len(l) }
func (l Int64List) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l Int64List) Less(i, j int) bool { return l[i] < l[j] }

// XYDataMapper x,y坐标系的数据集，用于图表数据
type XYDataMapper struct {
	Desc      bool                                `description:"表示倒序"`
	Round     int32                               `description:"保留几位小数"`
	Cache     map[string]*XYData                  `description:"表示过程中的缓存区"`
	List      XYDataList                          `description:"表示排序后的最终结果"`
	HandleY   func(data *XYData) (float64, error) `description:"自定义处理y的函数，data是同一个key的"`
	HandleAdd func(data *XYData, flags ...any)    `description:"自定义Add函数，也就是y的增减逻辑"`
}

func (th *XYDataMapper) Add(key string, x string, flags ...any) {
	if len(flags) == 0 {
		panic(exception.New("flags nil"))
	}
	if th.Cache == nil {
		th.Cache = map[string]*XYData{}
	}
	if _, ok := th.Cache[key]; !ok {
		th.Cache[key] = &XYData{X: x}
	}
	// 默认按递增处理
	if th.HandleAdd == nil {
		th.Cache[key].Y += cast.ToFloat64(flags[0])
	} else {
		th.HandleAdd(th.Cache[key], flags...)
	}
}

func (th *XYDataMapper) Result() XYDataList {
	if th.Cache != nil {
		arr := make(XYDataList, 0, len(th.Cache))
		for _, e := range th.Cache {
			if th.HandleY != nil {
				var err error
				e.Y, err = th.HandleY(e)
				if err != nil {
					// 存在错误则不加入list
					continue
				}
			}
			e.Desc = th.Desc
			if th.Round > 0 {
				e.Y = NewDecimal(e.Y).Round(th.Round).Float64()
			}
			arr = append(arr, e)
		}
		sort.Sort(arr)
		th.List = arr
	}
	return th.List
}

type XYData struct {
	X      string         `json:"x"`
	Y      float64        `json:"y"`
	Desc   bool           `json:"-" description:"表示是否倒序，用于Less函数"`
	Extend map[string]any `json:"extend"`
}
type XYDataList []*XYData

func (l XYDataList) Len() int      { return len(l) }
func (l XYDataList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l XYDataList) Less(i, j int) bool {
	if l[i].Desc {
		return cast.ToFloat64(l[i].Y) > cast.ToFloat64(l[j].Y)
	} else {
		return cast.ToFloat64(l[i].Y) < cast.ToFloat64(l[j].Y)
	}
}
