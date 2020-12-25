package timekit

import (
	"fmt"
	"log"
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
	log.Println(title + ": " + t.Start.Format(TimeLayoutAll))
	return t
}

func (th *TimeCost) PrintCost(msg string) {
	log.Print(fmt.Sprintf("%s: %dms", msg, time.Since(th.Start).Milliseconds()))
	th.Start = time.Now()
}
