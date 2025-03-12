package pq

import (
	"container/heap"
)

type PriorityQueueElement struct {
	value    any
	index    int
	priority float64
}

type PriorityQueue []*PriorityQueueElement

func (pq PriorityQueue) Len() int { 
  return len(pq) 
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*PriorityQueueElement)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil 
	*pq = old[:n-1]
	return item
}

func (pq *PriorityQueue) PushItem(value any, priority float64) {
	item := &PriorityQueueElement{value: value, priority: priority}
	heap.Push(pq, item)
}

func (pq *PriorityQueue) PopItem() (any, float64, bool) {
	if pq.Len() == 0 {
		var zero any
		return zero, 0, false
	}
	item := heap.Pop(pq).(*PriorityQueueElement)
	return item.value, item.priority, true
}
