package ddlib

import (
	"sync"
)

//Node .
type node struct {
	value interface{}
	next  *node
}

//Queue .
type Queue struct {
	size uint32
	head *node
	tail *node
	lock sync.Mutex
}

var nodePool = sync.Pool{
	New: func() interface{} {
		node := &node{}
		return node
	},
}

// NewQueue .
func NewQueue() *Queue {
	obj := &Queue{}
	obj.size = 0
	obj.head = nil
	obj.tail = nil
	return obj
}

//Pop .
func (q *Queue) Pop() interface{} {
	defer q.lock.Unlock()
	q.lock.Lock()
	if q.size > 0 {
		q.size--
		v := q.head.value
		tmp := q.head
		q.head = q.head.next
		tmp.value = nil
		tmp.next = nil
		nodePool.Put(tmp)
		return v
	}
	return nil
}

//Push .
func (q *Queue) Push(v interface{}) {
	defer q.lock.Unlock()
	q.lock.Lock()

	if q.size < CACHEMAX {
		n := nodePool.Get().(*node)
		n.value = v
		n.next = nil
		if nil != q.tail {
			q.tail.next = n
			q.tail = n
		} else {
			q.head = n
			q.tail = n
		}
		q.size++
	}
}
