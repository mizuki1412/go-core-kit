package concurrentkit

import (
	"github.com/mizuki1412/go-core-kit/v2/library/c"
	"sync"
)

type Group struct {
	wg sync.WaitGroup
}

func NewGroup() *Group {
	return &Group{}
}

func (g *Group) Add(f func(), shouldPanic bool) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if shouldPanic {
			f()
		} else {
			_ = c.RecoverFuncWrapper(f)
		}
	}()
}

func (g *Group) Process() {
	g.wg.Wait()
}
