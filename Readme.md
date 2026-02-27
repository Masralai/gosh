## GoSh (Go-Shell)

A custom, interactive command-line interface built in Go, designed for modularity and secure command execution.

---

### Key Features

* **Custom Command Engine:** Built on `urfave/cli/v3` for robust argument parsing and flag management.
* **Safety-First Logic:** Custom implementations for file operations with built-in recursive guards and confirmation prompts.
* **Isolated Environment:** Fully containerized setup via Docker for portable, zero-install execution.
* **Persistent History:** Integrated command scanner for a fluid, interactive terminal experience.

---

### Getting Started

#### Prerequisites

* **Option A:** [Go](https://go.dev/doc/install) (v1.24+)
* **Option B:** [Docker Desktop](https://www.docker.com/products/docker-desktop/)

#### Docker Execution (Recommended)

The easiest way to run GoSh without polluting your local environment:

```powershell
docker compose run gosh

```

*Note: This automatically builds the image and attaches an interactive TTY session.*

#### Manual Local Build

If you prefer to run natively on your machine:

1. Initialize: `go mod init go-cli`
2. Dependencies: `go get github.com/urfave/cli/v3`
3. Launch: `go run cli.go`
