package metrics

type Histogram struct {
	bucketSize int   // e.g. 10ms
	maxValue   int   // max latency we track (e.g. 10s = 10000ms)
	buckets    []int // counts
	total      int
}

func NewHistogram(bucketSize int, maxValue int) *Histogram {
	bucketCount := maxValue/bucketSize + 1

	return &Histogram{
		bucketSize: bucketSize,
		maxValue:   maxValue,
		buckets:    make([]int, bucketCount),
	}
}

func (h *Histogram) Add(latency int) {
	if latency > h.maxValue {
		latency = h.maxValue
	}
	index := latency / h.bucketSize
	h.buckets[index]++
	h.total++
}

func (h *Histogram) Percentile(p float64) int {
	if h.total == 0 {
		return 0
	}

	target := int(float64(h.total) * p / 100.0)
	count := 0

	for i, c := range h.buckets {
		count += c
		if count >= target {
			return i * h.bucketSize
		}
	}
	return h.maxValue
}
