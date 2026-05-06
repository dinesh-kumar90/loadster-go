package runner

import (
	"time"
)

type Pacer struct {
	ticker *time.Ticker
	tokens <-chan time.Time
}

func newPacer(rps int) *Pacer {
	if rps <= 0 {
		return &Pacer{}
	}

	interval := time.Second / time.Duration(rps)
	if interval <= 0 {
		interval = time.Nanosecond
	}
	ticker := time.NewTicker(interval)
	return &Pacer{ticker: ticker, tokens: ticker.C}
}

func (p *Pacer) Wait(end time.Time) bool {
	if p == nil || p.tokens == nil {
		return time.Now().Before(end)
	}

	remaining := time.Until(end)
	if remaining <= 0 {
		return false
	}

	select {
	case <-p.tokens:
		return time.Now().Before(end)
	case <-time.After(remaining):
		return false
	}
}

func (p *Pacer) Stop() {
	if p != nil && p.ticker != nil {
		p.ticker.Stop()
	}
}
