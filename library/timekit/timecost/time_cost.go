package timecost

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"time"
)

// TimeCost 程序执行耗时记录打点
type TimeCost struct {
	Start time.Time
}

type Params struct {
	showTitleMsg bool
}

func NewTimeCost(title string, params ...Params) TimeCost {
	t := TimeCost{
		Start: time.Now(),
	}
	param := Params{}
	if len(params) > 0 {
		param = params[0]
	}
	if param.showTitleMsg {
		logkit.Debug(title + ": " + t.Start.Format(timekit.TimeLayout))
	}
	return t
}

func (th *TimeCost) PrintCost(msg string) {
	logkit.Debug(fmt.Sprintf("%s: %dms", msg, time.Since(th.Start).Milliseconds()))
	th.Start = time.Now()
}
