package internal

import (
	"container/heap"
)

type EventQueue struct {
	_queue *InternalEventQueue
}

func CreateEventQueue() EventQueue {
	internalEventQueue := make(InternalEventQueue, 0)
	queue := EventQueue{
		_queue: &internalEventQueue,
	}
	heap.Init(queue._queue)
	return queue
}

func (queue EventQueue) Push(event *Event) {
	heap.Push(queue._queue, event)
}

func (queue EventQueue) Pop() *Event {
	return heap.Pop(queue._queue).(*Event)
}

func (queue EventQueue) Peek() *Event {
	return queue._queue.Peek().(*Event)
}

func (queue EventQueue) Len() int {
	return queue._queue.Len()
}

type InternalEventQueue []*Event

func (queue InternalEventQueue) Len() int {
	return len(queue)
}

func (queue InternalEventQueue) Less(i, j int) bool {
	return queue[i].DispatchAt < queue[j].DispatchAt
}

func (queue InternalEventQueue) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
	queue[i].Index = i
	queue[j].Index = j
}

func (queue *InternalEventQueue) Push(element any) {
	_queue := *queue

	event := element.(*Event)
	event.Index = len(_queue)

	*queue = append(_queue, event)
}

func (queue *InternalEventQueue) Pop() any {
	_queue := *queue
	_lastIndex := len(_queue) - 1

	item := _queue[_lastIndex]
	item.Index = -1

	_queue[_lastIndex] = nil
	*queue = _queue[:_lastIndex]

	return item
}

func (queue *InternalEventQueue) Peek() any {
	return (*queue)[0]
}
