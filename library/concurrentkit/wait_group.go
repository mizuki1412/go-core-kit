package concurrentkit

import (
	"sync"

	"github.com/mizuki1412/go-core-kit/v2/library/c"
)

type Group struct {
	wg sync.WaitGroup
}

func NewGroup() *Group {
	return &Group{}
}

func (g *Group) Add(f func(), shouldPanic bool) {
	g.wg.Go(func() {
		if shouldPanic {
			f()
		} else {
			_ = c.RecoverFuncWrapper(f)
		}
	})
}

func (g *Group) Process() {
	g.wg.Wait()
}
