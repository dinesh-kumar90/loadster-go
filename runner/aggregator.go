package runner

import (
	"fmt"
	"time"

	"github.com/dinesh-kumar90/loadster-go/metrics"
)

func aggregate(results <-chan Result, duration time.Duration) {
	stats := &metrics.Stats{}
	hist := metrics.NewHistogram(10, 10000) // 10ms buckets, up to 10s

	start := time.Now()

	for r := range results {
		lat := int(r.Latency.Milliseconds())

		stats.Add(float64(lat), r.Success)
		hist.Add(lat)
	}

	elapsed := time.Since(start).Seconds()

	fmt.Println("\n--- Load Test Results ---")
	fmt.Println("Total Requests:", stats.Total())
	fmt.Println("Success:", stats.Success())
	fmt.Println("Failed:", stats.Failed())

	fmt.Printf("RPS: %.2f\n", float64(stats.Total())/elapsed)

	fmt.Printf("\nLatency:\n")
	fmt.Printf("  Avg: %.2f ms\n", stats.Avg())
	fmt.Printf("  p50: %d ms\n", hist.Percentile(50))
	fmt.Printf("  p95: %d ms\n", hist.Percentile(95))
	fmt.Printf("  p99: %d ms\n", hist.Percentile(99))
}
