package go_common

import "sync"

type ArrayBlockQueue[E interface{}] interface {
	Put(e E)      // 阻塞放入
	Take() E      // 阻塞取出
	Add(e E) bool // 非阻塞放入
	Poll() E      // 非阻塞取出
	Peek() E      // 非阻塞查看
	IsEmpty() bool
	IsFull() bool
	Size() int
	Drain(max int) []E // 清空队列
}

func NewArrayBlockQueue[E interface{}](size int) ArrayBlockQueue[E] {
	return &arrayBlockQueue[E]{
		elements:  make([]E, size),
		size:      size,
		count:     0,
		putIndex:  0,
		takeIndex: 0,
		cond:      sync.NewCond(&sync.Mutex{}),
	}
}

type arrayBlockQueue[E interface{}] struct {
	elements  []E
	size      int
	count     int
	putIndex  int
	takeIndex int
	cond      *sync.Cond
}

func (q *arrayBlockQueue[E]) Put(e E) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for q.size == q.count {
		q.cond.Wait()
	}
	q.enqueue(e)
}

func (q *arrayBlockQueue[E]) Take() E {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for q.count == 0 {
		q.cond.Wait()
	}

	return q.dequeue()
}

func (q *arrayBlockQueue[E]) Add(e E) bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.count == q.size {
		return false
	}
	q.enqueue(e)
	return true
}

func (q *arrayBlockQueue[E]) Poll() (e E) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.count == 0 {
		return
	}
	return q.dequeue()
}

func (q *arrayBlockQueue[E]) Peek() (e E) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.count == 0 {
		return
	}
	e = q.elements[q.takeIndex]
	return e
}

func (q *arrayBlockQueue[E]) IsEmpty() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.count == 0
}

func (q *arrayBlockQueue[E]) IsFull() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.count == q.size
}

func (q *arrayBlockQueue[E]) Size() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.count
}

func (q *arrayBlockQueue[E]) Drain(max int) []E {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.count == 0 {
		return nil
	}
	if max > q.count {
		max = q.count
	}
	drain := make([]E, max)
	for i := 0; i < max; i++ {
		drain[i] = q.dequeue()
	}
	return drain
}

func (q *arrayBlockQueue[E]) clear(index int) {
	var e E
	q.elements[index] = e
}

func (q *arrayBlockQueue[E]) enqueue(e E) {
	q.elements[q.putIndex] = e
	q.putIndex = (q.putIndex + 1) % q.size
	q.count++
	q.cond.Broadcast()
}
func (q *arrayBlockQueue[E]) dequeue() E {
	e := q.elements[q.takeIndex]
	q.clear(q.takeIndex)
	q.takeIndex = (q.takeIndex + 1) % q.size
	q.count--
	q.cond.Broadcast()
	return e
}
