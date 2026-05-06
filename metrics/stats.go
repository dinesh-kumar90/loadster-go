package metrics

type Stats struct {
	total   int
	success int
	failed  int
	sum     float64
}

func (s *Stats) Add(latency float64, success bool) {
	s.total++
	s.sum += latency

	if success {
		s.success++
	} else {
		s.failed++
	}
}

func (s *Stats) Total() int {
	return s.total
}

func (s *Stats) Success() int {
	return s.success
}

func (s *Stats) Failed() int {
	return s.failed
}

func (s *Stats) Avg() float64 {
	if s.total == 0 {
		return 0
	}
	return s.sum / float64(s.total)
}
