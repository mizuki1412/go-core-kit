package class

import "sync"

// 用于异步等待

type LockWait struct {
	sync.RWMutex
	Data map[string]chan map[string]any
}

func (th *LockWait) Open(cid string) {
	th.Lock()
	defer th.Unlock()
	th.Data[cid] = make(chan map[string]any)
}

func (th *LockWait) Close(cid string) {
	th.Lock()
	defer th.Unlock()
	delete(th.Data, cid)
}

func (th *LockWait) ChanExist(cid string) bool {
	_, ok := th.Data[cid]
	return ok
}

func (th *LockWait) Chan(cid string) chan map[string]any {
	return th.Data[cid]
}
