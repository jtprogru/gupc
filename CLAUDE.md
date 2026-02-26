# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common development commands

The project uses a **Taskfile.yml** for most routine operations. Below are the most frequently used tasks (run with `task <task-name>`):

| Task | Description |
|------|-------------|
| `run:cmd` | Run the program via `go run main.go <SLA‑percent> [--json]` (use `{{ .CLI_ARGS }}` to pass arguments). |
| `run:bin` | Execute the compiled binary located at `./dist/gupc`. Requires the binary to be built first. |
| `build:bin` | Build a statically linked binary (`CGO_ENABLED=0 go build -o ./dist/gupc main.go`). |
| `tidy` | Ensure module dependencies are up‑to‑date (`go mod tidy`). |
| `fmt` | Reformat the source tree (`gofmt -s -w .`). |
| `vet` | Run Go vet (`go vet ./...`). |
| `lint` | Run `golangci-lint -v run`. |
| `test` | Run the full test suite (`task test:short`). |
| `test:short` | Run short tests with coverage (`go test --short -coverprofile=cover.out -v ./...`). |
| `test:coverage` | Run all tests with coverage collection (`go test -coverprofile=cover.out -v ./...`). |
| `test:race` | Run tests with the race detector (`go test -race -coverprofile=cover.out -v ./...`). |
| `test:bench` | Run benchmarks (`go test -v -bench=. -benchmem -run=^$ ./...`). |
| `test:watch` | Continuously run tests on file changes (requires `watchexec`). |
| `install:global` | Install the module globally (`go install`). |
| `install:dev` | Build and install the binary into `$GOPATH/bin`. |

### Running the CLI directly

```bash
# Quick run without building:
$ go run main.go 99.9            # human‑readable output
$ go run main.go 99.9 --json     # JSON output
```

```bash
# After building the binary:
$ task build:bin
$ ./dist/gupc 99.9               # human‑readable output
$ ./dist/gupc 99.9 --json        # JSON output
```

## High‑level architecture

- **Entry point**: `main.go` imports the `cmd` package and invokes `cmd.Execute()`.
- **Command definition**: `cmd/root.go` defines a Cobra root command (`Use: "gupc"`). It parses a single positional argument – the SLA percentage – and a `--json` flag.
- **Core logic**:
  1. Parse the SLA argument as a `float64` percentage.
  2. Validate the value (must be > 0 and < 100).
  3. Compute the number of nines (e.g., "three nines" for 99.9 %).
  4. Derive downtime for various periods (day, week, month, quarter, year) based on the SLA.
  5. If `--json` is set, marshal a `SLAInfo` struct to formatted JSON; otherwise, print a human‑readable report using `formatDuration`.
- **Utility**: `formatDuration` converts a duration in seconds to a compact string like `1h23m45s`.
- **Project layout**:
  - `cmd/` – contains the Cobra command implementation.
  - `main.go` – thin wrapper that just calls `cmd.Execute()`.
  - `go.mod` – module definition (`github.com/jtprogru/gupc`) and declares the only external dependency: `github.com/spf13/cobra`.
- **Testing & linting**: No tests are currently present, but the Taskfile includes standard Go test, vet, and lint commands for future development.

## Notable configuration files

- `go.mod` – declares the module path and Go version (`go 1.25`).
- `Taskfile.yml` – provides a high‑level workflow (build, test, lint, etc.).
- `plan.md` – contains an unfinished development plan for improving input validation and output formatting.

---

*This file is intended for Claude Code to quickly understand how to build, run, and work with the repository.*