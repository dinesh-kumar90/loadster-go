# Loadster Go

A lightweight, high-concurrency CLI load testing tool written in Go.

Loadster Go is designed to stress-test HTTP endpoints with concurrent workers, optional global RPS pacing, and startup ramping to avoid sudden thread spikes.

## Why This Project

I built this project to explore practical performance engineering in Go:
- Concurrency orchestration with goroutines + channels
- Shared HTTP transport tuning for high-throughput workloads
- Runtime metrics aggregation (RPS, success/failure counts, latency percentiles)
- CLI ergonomics for repeatable load-testing experiments

## Features

- Concurrent virtual users (`-users`)
- Time-bound test runs (`-duration`)
- HTTP method selection (`-method`)
- Custom request body support (`-body`)
- Custom request headers (`-header "Key: Value"`; repeatable)
- Optional global rate limiting (`-rps`)
- Worker ramp-up controls (`-ramp-step`, `-ramp-interval`)
- Shared and tuned HTTP client/transport for better connection reuse
- Summary output with:
  - Total requests
  - Success/failed requests
  - Effective RPS
  - Latency Avg / p50 / p95 / p99

## Tech Stack

- Go 1.26+
- Standard library only (`flag`, `net/http`, goroutines/channels)

## Project Structure

```text
.
├── main.go                # CLI flags and entrypoint
├── runner/
│   ├── runner.go          # test orchestration + ramp-up + client setup
│   ├── worker.go          # request execution workers
│   ├── aggregator.go      # result aggregation and reporting
│   └── pacer.go           # global RPS pacing
├── metrics/
│   ├── stats.go           # totals, success/failure, avg latency
│   └── histogram.go       # percentile calculations
└── README.md
```

## Installation

```bash
git clone https://github.com/dinesh-kumar90/loadster-go.git
cd loadster-go
go build -o loadster .
```

## Binary Downloads

Prebuilt binaries are available on the GitHub Releases page:

- https://github.com/dinesh-kumar90/loadster-go/releases

Download the asset matching your OS/architecture and run it directly.
Releases are built automatically by GitHub Actions when you push a tag like `v1.0.0`.

## Usage

```bash
./loadster -url=<TARGET_URL> [flags]
```

### Flags

- `-url` string (required): Target URL
- `-users` int (default: `10`): Number of concurrent workers
- `-duration` duration (default: `10s`): Total test duration
- `-method` string (default: `GET`): HTTP method
- `-body` string (default: empty): Request body payload
- `-header` string (repeatable): Custom header in `Key: Value` format
- `-rps` int (default: `0`): Target global requests/second (`0` = unthrottled)
- `-ramp-step` int (default: `100`): Workers started per ramp batch
- `-ramp-interval` duration (default: `50ms`): Delay between ramp batches

## Example Commands

Unthrottled concurrency test:

```bash
./loadster -url=http://localhost:8080 -users=3000 -duration=30s
```

RPS-controlled test:

```bash
./loadster -url=http://localhost:8080 -users=2000 -duration=45s -rps=10000
```

POST with JSON body and headers:

```bash
./loadster \
  -url=http://localhost:8080/api/orders \
  -method=POST \
  -body='{\"item\":\"book\",\"qty\":2}' \
  -header='Content-Type: application/json' \
  -header='Authorization: Bearer test-token' \
  -users=500 \
  -duration=30s \
  -rps=2000
```

Large-user test with gradual ramp-up:

```bash
./loadster -url=http://localhost:8080 -users=20000 -duration=60s -rps=15000 -ramp-step=250 -ramp-interval=100ms
```

## Sample Output

```text
Starting load test...
URL: http://localhost:8080
Users: 2000
Duration: 45s
RPS: 10000
Ramp Step: 100
Ramp Interval: 50ms

--- Load Test Results ---
Total Requests: 450012
Success: 448903
Failed: 1109
RPS: 10000.26

Latency:
  Avg: 21.37 ms
  p50: 18 ms
  p95: 54 ms
  p99: 88 ms
```

## Notes

- This tool is intended for performance testing environments and systems you have permission to test.
- Very high concurrency may require OS/container tuning (for example thread/process and file-descriptor limits).

## Roadmap

- Status-code distribution reporting
- JSON/CSV output mode
- Distributed load generation mode

## Create Release Binaries

Use the included script to build multi-OS binaries and checksums:

```bash
./scripts/release.sh v1.0.0
```

Or use CI/CD: push a version tag and GitHub Actions will build and attach all binaries automatically:

```bash
git tag v1.0.0
git push origin v1.0.0
```

This generates artifacts in `dist/`:
- `loadster-v1.0.0-linux-amd64`
- `loadster-v1.0.0-linux-arm64`
- `loadster-v1.0.0-darwin-amd64`
- `loadster-v1.0.0-darwin-arm64`
- `loadster-v1.0.0-windows-amd64.exe`
- `loadster-v1.0.0-windows-arm64.exe`
- `checksums.txt`

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
