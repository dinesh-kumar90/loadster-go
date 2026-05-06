package runner

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	URL          string
	Users        int
	Duration     time.Duration
	Method       string
	Body         string
	Headers      http.Header
	RPS          int
	RampStep     int
	RampInterval time.Duration
}

func Run(cfg Config) {
	if cfg.Users <= 0 {
		return
	}
	if cfg.RampStep <= 0 {
		cfg.RampStep = cfg.Users
	}
	if cfg.RampInterval < 0 {
		cfg.RampInterval = 0
	}

	results := make(chan Result, 10000)

	var wg sync.WaitGroup
	end := time.Now().Add(cfg.Duration)
	client := newHTTPClient(cfg.Users)
	pacer := newPacer(cfg.RPS)
	defer pacer.Stop()

	started := 0
	for started < cfg.Users {
		batch := cfg.RampStep
		remaining := cfg.Users - started
		if batch > remaining {
			batch = remaining
		}

		for i := 0; i < batch; i++ {
			wg.Add(1)
			go Worker(cfg.URL, cfg.Method, cfg.Body, cfg.Headers, end, client, pacer, results, &wg)
		}
		started += batch

		if started < cfg.Users && cfg.RampInterval > 0 {
			time.Sleep(cfg.RampInterval)
		}
	}

	// Close channel when workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate results
	aggregate(results, cfg.Duration)
}

func newHTTPClient(users int) *http.Client {
	maxConns := users * 4
	if maxConns < 1000 {
		maxConns = 1000
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          maxConns,
		MaxIdleConnsPerHost:   maxConns,
		MaxConnsPerHost:       maxConns,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &http.Client{
		Timeout:   15 * time.Second,
		Transport: transport,
	}
}
