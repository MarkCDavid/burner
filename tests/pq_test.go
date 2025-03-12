package tests

import (
	"container/heap"
	"testing"

	pkgPriorityQueue "github.com/MarkCDavid/burner/internal/pq"
)

func TestPriorityQueue(t *testing.T) {
	pq := &pkgPriorityQueue.PriorityQueue{}
	heap.Init(pq)

	pq.PushItem(1, 3.0)
	pq.PushItem(2, 1.0)
	pq.PushItem(3, 2.0)

	if pq.Len() != 3 {
		t.Errorf("Expected queue length 3, got %d", pq.Len())
	}

	val, priority, ok := pq.PopItem()
	if !ok || val != 2 || priority != 1.0 {
		t.Errorf("Expected (2, 1.0), got (%v, %v)", val, priority)
	}

	val, priority, ok = pq.PopItem()
	if !ok || val != 3 || priority != 2.0 {
		t.Errorf("Expected (3, 2.0), got (%v, %v)", val, priority)
	}

	val, priority, ok = pq.PopItem()
	if !ok || val != 1 || priority != 3.0 {
		t.Errorf("Expected (1, 3.0), got (%v, %v)", val, priority)
	}
}

func TestPopFromEmptyQueue(t *testing.T) {
	pq := &pkgPriorityQueue.PriorityQueue{}
	heap.Init(pq)

	val, priority, ok := pq.PopItem()
	if ok {
		t.Errorf("Expected empty queue pop to fail, but got (%v, %v)", val, priority)
	}
}

func TestPriorityQueueStrings(t *testing.T) {
	pq := &pkgPriorityQueue.PriorityQueue{}
	heap.Init(pq)

	pq.PushItem("low", 5.0)
	pq.PushItem("medium", 3.0)
	pq.PushItem("high", 1.0)

	val, _, _ := pq.PopItem()
	if val != "high" {
		t.Errorf("Expected 'high' priority item, got %v", val)
	}

	val, _, _ = pq.PopItem()
	if val != "medium" {
		t.Errorf("Expected 'medium' priority item, got %v", val)
	}

	val, _, _ = pq.PopItem()
	if val != "low" {
		t.Errorf("Expected 'low' priority item, got %v", val)
	}
}

func TestPriorityQueueSwap(t *testing.T) {
	pq := &pkgPriorityQueue.PriorityQueue{}
	heap.Init(pq)

	pq.PushItem(1, 3.0)
	pq.PushItem(2, 1.0)
	pq.PushItem(3, 2.0)

	heap.Fix(pq, 0) // Force a reordering

	if pq.Len() != 3 {
		t.Errorf("Expected queue length 3, got %d", pq.Len())
	}
}
