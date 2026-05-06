package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/dinesh-kumar90/loadster-go/runner"
)

func main() {
	url := flag.String("url", "", "Target URL")                            // -url=http://example.com
	users := flag.Int("users", 10, "Concurrent users")                     // -users=10000 default=10
	duration := flag.Duration("duration", 10*time.Second, "Test duration") // -duration=10s default=10s
	method := flag.String("method", "GET", "HTTP method")                  // -method=GET default=GET
	rps := flag.Int("rps", 0, "Target requests per second (0 disables rate limiting)")
	rampStep := flag.Int("ramp-step", 100, "Workers to start per ramp step")
	rampInterval := flag.Duration("ramp-interval", 50*time.Millisecond, "Delay between ramp steps")

	flag.Parse()
	if *url == "" {
		panic("URL is required")
	}
	fmt.Println("Starting load test...")
	fmt.Println("URL:", *url)
	fmt.Println("Users:", *users)
	fmt.Println("Duration:", *duration)
	if *rps > 0 {
		fmt.Println("RPS:", *rps)
	}
	fmt.Println("Ramp Step:", *rampStep)
	fmt.Println("Ramp Interval:", *rampInterval)

	cfg := runner.Config{
		URL:          *url,
		Users:        *users,
		Duration:     *duration,
		Method:       *method,
		RPS:          *rps,
		RampStep:     *rampStep,
		RampInterval: *rampInterval,
	}

	runner.Run(cfg)
}
