package runner

import (
	"net/http"
	"sync"
	"time"
)

type Result struct {
	Latency time.Duration
	Success bool
}

func Worker(url string, method string, end time.Time, client *http.Client, pacer *Pacer, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for time.Now().Before(end) {
		if !pacer.Wait(end) {
			return
		}

		start := time.Now()

		req, _ := http.NewRequest(method, url, nil)
		resp, err := client.Do(req)

		latency := time.Since(start)

		if err != nil {
			results <- Result{Latency: latency, Success: false}
			continue
		}

		resp.Body.Close()

		results <- Result{
			Latency: latency,
			Success: resp.StatusCode < 500,
		}
	}
}
