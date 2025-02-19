package class

import "sync"

type (
	//Queue 队列
	Queue[T any] struct {
		top    *node[T]
		rear   *node[T]
		length int
		lock1  sync.Mutex
		lock2  sync.Mutex
	}
	//双向链表节点
	node[T any] struct {
		pre   *node[T]
		next  *node[T]
		value T
	}
)

// NewQueue Create a new queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Len 获取队列长度
func (th *Queue[T]) Len() int {
	return th.length
}

// Peek 返回队列顶端元素
func (th *Queue[T]) Peek() *T {
	if th.top == nil {
		return nil
	}
	return &th.top.value
}

// Push 入队操作
func (th *Queue[T]) Push(v T) {
	th.lock1.Lock()
	defer th.lock1.Unlock()
	n := &node[T]{nil, nil, v}
	if th.length == 0 {
		th.top = n
		th.rear = th.top
	} else {
		n.pre = th.rear
		th.rear.next = n
		th.rear = n
	}
	th.length++
}

// Pop 出队操作
func (th *Queue[T]) Pop() *T {
	th.lock1.Lock()
	defer th.lock1.Unlock()
	if th.length == 0 {
		return nil
	}
	n := th.top
	if th.top.next == nil {
		th.top = nil
	} else {
		th.top = th.top.next
		th.top.pre.next = nil
		th.top.pre = nil
	}
	th.length--
	return &n.value
}
