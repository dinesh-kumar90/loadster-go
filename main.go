package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
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
	body := flag.String("body", "", "Request body payload")
	var headers headerFlags
	flag.Var(&headers, "header", "Custom request header in 'Key: Value' format (repeatable)")

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
	if *body != "" {
		fmt.Println("Body: provided")
	}
	if len(headers) > 0 {
		fmt.Println("Headers:", len(headers))
	}

	cfg := runner.Config{
		URL:          *url,
		Users:        *users,
		Duration:     *duration,
		Method:       *method,
		RPS:          *rps,
		RampStep:     *rampStep,
		RampInterval: *rampInterval,
		Body:         *body,
		Headers:      parseHeaders(headers),
	}

	runner.Run(cfg)
}

type headerFlags []string

func (h *headerFlags) String() string {
	return strings.Join(*h, ",")
}

func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func parseHeaders(entries []string) http.Header {
	out := make(http.Header)
	for _, entry := range entries {
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}
		out.Add(key, value)
	}
	return out
}
