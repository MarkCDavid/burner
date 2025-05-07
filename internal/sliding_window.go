package internal

type SWFloat64 struct {
	size  int64
	start int64
	count int64

	values []float64
	total  float64
}

func NewSWFloat64(size int64) *SWFloat64 {
	return &SWFloat64{
		values: make([]float64, size),
		size:   size,
	}
}

func (sw *SWFloat64) Add(value float64) {
	if sw.count == sw.size {
		sw.total -= sw.values[sw.start]
	} else {
		sw.count++
	}

	sw.values[sw.start] = value
	sw.total += value
	sw.start = (sw.start + 1) % sw.size
}

func (sw *SWFloat64) Sum() float64 {
	return sw.total
}

func (sw *SWFloat64) Average() float64 {
	if sw.count == 0 {
		return 0
	}
	return sw.total / float64(sw.count)
}

type SWCounterInt64[T comparable] struct {
	size  int64
	start int64
	count int64

	values []T
	counts map[T]int64
}

func NewSWCounterInt64[T comparable](size int64) *SWCounterInt64[T] {
	return &SWCounterInt64[T]{
		values: make([]T, size),
		counts: make(map[T]int64),
		size:   size,
	}
}

func (sw *SWCounterInt64[T]) Add(value T) {
	if sw.count == sw.size {
		sw.counts[value]--
	} else {
		sw.count++
	}

	sw.values[sw.start] = value
	sw.counts[value]++
	sw.start = (sw.start + 1) % sw.size
}

func (sw *SWCounterInt64[T]) Get(value T) int64 {
	return sw.counts[value]
}
