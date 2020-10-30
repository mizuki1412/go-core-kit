package class

import "sync"

type (
	//Queue 队列
	Queue struct {
		top    *node
		rear   *node
		length int
		sync.RWMutex
	}
	//双向链表节点
	node struct {
		pre   *node
		next  *node
		value interface{}
	}
)

// Create a new queue
//func New() *Queue {
//	return &Queue{}
//}
//获取队列长度
func (th *Queue) Len() int {
	return th.length
}

//返回true队列不为空
func (th *Queue) Any() bool {
	return th.length > 0
}

//返回队列顶端元素
func (th *Queue) Peek() interface{} {
	if th.top == nil {
		return nil
	}
	return th.top.value
}

//入队操作
func (th *Queue) Push(v interface{}) {
	th.Lock()
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
	th.Unlock()
}

//出队操作
func (th *Queue) Pop() interface{} {
	th.Lock()
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
	th.Unlock()
	return n.value
}
