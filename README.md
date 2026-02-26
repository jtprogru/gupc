# gupc – Uptime & Downtime Calculator

**gupc** is a small command‑line utility that converts an SLA percentage (e.g., `99.9`) into:

* The human‑readable “number of nines” (e.g., *three nines* for `99.9%`).
* Expected downtime for a day, week, month, quarter, and year.

It can output the result as plain text or as a JSON object.

---

## Table of Contents
1. [Overview](#overview)
2. [Installation](#installation)
   - [From source (recommended)](#from-source)
   - [Using `go install`](#go-install)
3. [Running the CLI](#running-the-cli)
   - [Quick run (`go run`)](#quick-run)
   - [Run compiled binary](#run-binary)
4. [Development workflow](#development-workflow)
   - [Taskfile commands](#taskfile-commands)
   - [Formatting & linting](#formatting-and-linting)
   - [Testing](#testing)
5. [Project structure](#project-structure)
6. [Contributing](#contributing)
7. [License](#license)

---

## Overview <a name="overview"></a>

The program accepts a single positional argument – the SLA **percentage** – and an optional `--json` flag.

* If `--json` is provided, the tool prints a pretty‑printed JSON representation of the `SLAInfo` struct.
* Otherwise it prints a concise, human‑readable report.

Example output (human‑readable, `99.9`):

```
99.90% SLA (three nines)
Daily downtime: 8m40s
Weekly downtime: 1h0m21s
Monthly downtime: 4h16m24s
Quarterly downtime: 12h49m12s
Yearly downtime: 43h18m26s
```

Example JSON (`99.9 --json`):

```json
{
    "SLA": 99.9,
    "Nines": "three nines",
    "DailyDownSecs": 520.0,
    "DailyDown": "8m40s",
    "WeeklyDownSecs": 3640.0,
    "WeeklyDown": "1h0m40s",
    "MonthlyDownSecs": 15600.0,
    "MonthlyDown": "4h20m0s",
    "QuarterlyDownSecs": 46800.0,
    "QuarterlyDown": "12h60m0s",
    "YearlyDownSecs": 157680.0,
    "YearlyDown": "43h48m0s"
}
```

---

## Installation <a name="installation"></a>

### From source (recommended) <a name="from-source"></a>

1. **Clone the repository**

   ```bash
   git clone https://github.com/jtprogru/gupc.git
   cd gupc
   ```

2. **Ensure you have Go 1.25+ installed** (the module declares `go 1.25`).

3. **Install dependencies**

   ```bash
   task tidy               # runs `go mod tidy`
   ```

4. **Build the binary**

   ```bash
   task build:bin          # creates ./dist/gupc
   ```

   The binary is statically linked (`CGO_ENABLED=0`) and ready to run on any Linux/macOS host.

### `go install` <a name="go-install"></a>

If you prefer a one‑liner:

```bash
go install github.com/jtprogru/gupc@latest
```

The binary will be placed in `$GOPATH/bin` (or `$HOME/go/bin`). Ensure that directory is on your `PATH`.

---

## Running the CLI <a name="running-the-cli"></a>

### Quick run (`go run`) <a name="quick-run"></a>

```bash
go run main.go <SLA-percent> [--json]
```

Example:

```bash
go run main.go 99.9
go run main.go 99.9 --json
```

### Run compiled binary <a name="run-binary"></a>

```bash
./dist/gupc <SLA-percent> [--json]
```

If you installed globally with `go install`:

```bash
gupc 99.9
gupc 99.9 --json
```

---

## Development workflow <a name="development-workflow"></a>

The project uses **Taskfile.yml** (powered by https://taskfile.dev) to encapsulate common commands.

| Task | Description |
|------|-------------|
| `run:cmd` | Run via `go run main.go …` (use `{{ .CLI_ARGS }}` to pass args). |
| `run:bin` | Execute the compiled binary (`./dist/gupc`). |
| `build:bin` | Build static binary (`CGO_ENABLED=0 go build -o ./dist/gupc main.go`). |
| `tidy` | `go mod tidy` – ensure module dependencies are tidy. |
| `fmt` | Run `gofmt -s -w .` to format source. |
| `vet` | Run `go vet ./...` for static analysis. |
| `lint` | Run `golangci-lint -v run`. |
| `test` | Run all tests (`task test:short`). |
| `test:short` | Run short tests with coverage (`go test --short -coverprofile=cover.out -v ./...`). |
| `test:coverage` | Run full test suite with coverage (`go test -coverprofile=cover.out -v ./...`). |
| `test:race` | Run tests with the race detector. |
| `test:bench` | Run benchmarks (`go test -bench=. -benchmem -run=^$ ./...`). |
| `test:watch` | Continuous testing on file changes (requires `watchexec`). |
| `install:global` | `go install` – install the module globally. |
| `install:dev` | Build binary and copy it into `$GOPATH/bin`. |

### Formatting & linting <a name="formatting-and-linting"></a>

```bash
task fmt
task vet
task lint
```

All three should pass before committing changes.

### Testing <a name="testing"></a>

There are currently no unit tests, but the scaffolding is ready:

```bash
task test:short      # fast tests with coverage
task test:coverage   # full suite with coverage report
task test:race       # race detection
task test:bench      # benchmarks
```

Add tests under `*_test.go` files as needed; `go test ./...` will automatically discover them.

---

## Project structure <a name="project-structure"></a>

```
.
├── CLAUDE.md          # Guidance for Claude Code (generated)
├── README.md          # **You are reading it**
├── Taskfile.yml       # Task automation definition
├── go.mod             # Module definition (go 1.25)
├── go.sum
├── main.go            # Application entry point; calls cmd.Execute()
├── cmd/
│   └── root.go       # Cobra command implementation, core logic
└── plan.md           # Development plan (future improvements)
```

* **`main.go`** – Tiny wrapper that delegates to the `cmd` package.
* **`cmd/root.go`** – Defines the Cobra root command (`gupc`). Handles:
  1. Argument parsing (`float64` SLA percent).
  2. Validation (`0 < percent < 100`).
  3. Calculation of “nines” and downtime periods.
  4. Output formatting (human‑readable vs JSON).
* **`formatDuration`** – Helper that turns seconds into `XhYmZs` strings.

---

## Contributing <a name="contributing"></a>

1. Fork the repository.
2. Create a feature or bug‑fix branch.
3. Ensure `task fmt && task vet && task lint && task test` all succeed.
4. Submit a pull request.

Feel free to open issues for:

* Feature requests (e.g., supporting additional output formats).
* Bugs in the SLA calculation or CLI parsing.
* Documentation improvements.

---

## License <a name="license"></a>

This project is released under the **MIT License** – see the `LICENSE` file for details.
