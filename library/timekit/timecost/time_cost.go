package timecost

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
	"time"
)

// TimeCost 程序执行耗时记录打点
type TimeCost struct {
	Title string
	Start time.Time
	//Unit string
}

type Params struct {
	//Unit string // 时间显示单位 s,m
}

func NewTimeCost(title string, params ...Params) TimeCost {
	t := TimeCost{
		Title: title,
		Start: time.Now(),
	}
	//param := Params{}
	//if len(params) > 0 {
	//	param = params[0]
	//}
	//t.Unit = param.Unit
	return t
}

// 每次打印后重新计时
func (th *TimeCost) PrintCost(msg string, logInfo ...bool) {
	message := fmt.Sprintf("time-cost【%s|%s】%fs", th.Title, msg, time.Since(th.Start).Seconds())
	if len(logInfo) > 0 && logInfo[0] {
		logkit.Info(message)
	} else {
		logkit.Debug(message)
	}
	th.Start = time.Now()
}
