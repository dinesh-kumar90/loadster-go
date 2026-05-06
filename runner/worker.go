package runner

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	Latency time.Duration
	Success bool
}

func Worker(url string, method string, body string, headers http.Header, end time.Time, client *http.Client, pacer *Pacer, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for time.Now().Before(end) {
		if !pacer.Wait(end) {
			return
		}

		start := time.Now()

		req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
		if err != nil {
			results <- Result{Latency: 0, Success: false}
			continue
		}
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
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
