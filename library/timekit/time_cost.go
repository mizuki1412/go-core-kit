package timekit

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"time"
)

// 程序执行耗时记录打点
type TimeCost struct {
	Start time.Time
}

func NewTimeCost(title string) TimeCost {
	t := TimeCost{
		Start: time.Now(),
	}
	logkit.Debug(title + ": " + t.Start.Format(TimeLayout))
	return t
}

func (th *TimeCost) PrintCost(msg string) {
	logkit.Debug(fmt.Sprintf("%s: %dms", msg, time.Since(th.Start).Milliseconds()))
	th.Start = time.Now()
}
