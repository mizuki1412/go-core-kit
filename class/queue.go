package class

import "sync"

type (
	//Queue 队列
	Queue struct {
		top    *node
		rear   *node
		length int
		lock1  sync.Mutex
		lock2  sync.Mutex
	}
	//双向链表节点
	node struct {
		pre   *node
		next  *node
		value interface{}
	}
)

// NewQueue Create a new queue
func NewQueue() *Queue {
	return &Queue{}
}

// Len 获取队列长度
func (th *Queue) Len() int {
	return th.length
}

// Peek 返回队列顶端元素
func (th *Queue) Peek() interface{} {
	if th.top == nil {
		return nil
	}
	return th.top.value
}

// Push 入队操作
func (th *Queue) Push(v interface{}) {
	th.lock1.Lock()
	defer th.lock1.Unlock()
	n := &node{nil, nil, v}
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
func (th *Queue) Pop() interface{} {
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
	return n.value
}
