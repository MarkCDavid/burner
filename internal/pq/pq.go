package pq

import (
	"container/heap"
)

type PriorityQueueElement[T any] struct {
	value    T
	index    int
	priority float64
}

type PriorityQueue[T any] []*PriorityQueueElement[T]

func (pq PriorityQueue[T]) Len() int { 
  return len(pq) 
}

func (pq PriorityQueue[T]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue[T]) Push(x any) {
	item := x.(*PriorityQueueElement[T])
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil 
	*pq = old[:n-1]
	return item
}

func (pq *PriorityQueue[T]) PushItem(value T, priority float64) {
	item := &PriorityQueueElement[T]{value: value, priority: priority}
	heap.Push(pq, item)
}

func (pq *PriorityQueue[T]) PopItem() (T, float64, bool) {
	if pq.Len() == 0 {
		var zero T
		return zero, 0, false
	}
	item := heap.Pop(pq).(*PriorityQueueElement[T])
	return item.value, item.priority, true
}
