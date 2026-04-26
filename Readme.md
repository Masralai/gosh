# GoSh (Go-Shell)

[![Go Report Card](https://goreportcard.com/badge/github.com/Masralai/gosh)](https://goreportcard.com/report/github.com/Masralai/gosh)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Masralai/gosh)](https://github.com/Masralai/gosh)

GoSh is a modular, interactive command-line interface built in Go. It provides a containerized environment for secure command execution, system monitoring, and file management.

---

## Table of Contents

- [Key Features](#key-features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Building](#building)
- [Testing](#testing)
- [Command Reference](#command-reference)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

---

## Key Features

- **Robust Command Engine:** Utilizes `urfave/cli/v3` for sophisticated argument parsing and flag management.
- **Enhanced Security:** Path traversal guards, directory traversal checks, and confirmation prompts for destructive actions.
- **Portability:** Works natively or containerized with Docker.
- **Interactive Experience:** Integrated command scanning provides a fluid terminal interaction.
- **System Monitoring:** Built-in tools for monitoring CPU, memory, disk usage, and process management.
- **Network Utilities:** Integrated network diagnostic tools including ping functionality.
- **Archive Management:** Support for creating and extracting ZIP archives with security checks.
- **Security Scanned:** Code undergoes gosec security analysis.

---

## Architecture

```
gosh/
├── cmd/gosh/main.go       # Entry point
├── internal/handlers/     # Command implementations
│   ├── commands.go       # Command aggregator
│   ├── shell.go        # Shell utilities (ls, cd, pwd, mkdir, cp, mv, etc.)
│   ├── system.go      # System monitoring (sys, mu, du, ps, kill, etc.)
│   ├── text.go      # Text processing (grep, head, tail)
│   ├── network.go   # Networking (ping)
│   └── storage.go   # Archives (zip, unzip)
├── Makefile          # Build automation
├── Dockerfile       # Container definition
└── compose.yaml    # Docker Compose
```

### Component Overview

| Component | Description |
|-----------|-------------|
| `cmd/gosh/main.go` | Entry point, sets up urfave/cli and command loop |
| `internal/handlers/` | All command implementations |

### Command Dispatcher

The core engine is built on `urfave/cli/v3`. It maps user input to specific handler functions in `internal/handlers/`. New commands can be added by implementing a handler function and registering it in `commands.go`.

---

## Prerequisites

- **Go:** Version 1.25 or higher
- **Docker:** Optional, for containerized execution

---

## Installation

### Native (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/Masralai/gosh
cd gosh
```

2. Build:
```bash
make build
# or
go build -o gosh ./cmd/gosh
```

3. Run:
```bash
./gosh
```

### Docker

```bash
docker compose up --build
```

---

## Building

### Using Makefile

```bash
make build    # Builds binary to ./gosh
make test    # Runs tests with race detector
make lint   # Runs linter
make vet    # Runs go vet
make clean  # Removes binary
```

### Manual Build

```bash
# Build binary
go build -o gosh ./cmd/gosh

# Run tests
go test -race ./...

# Lint
golangci-lint run ./...

# Vet
go vet ./...
```

---

## Testing

Run tests with race detection:

```bash
make test
# or
go test -race ./...
```

---

## Command Reference

### File Operations

| Command | Description | Flags |
| :--- | :--- | :--- |
| `ls` | List directory contents | `-a` (hidden), `-R` (recursive), `-S` (size) |
| `cd` | Change directory | |
| `pwd` | Print working directory | |
| `mkdir` | Create directory | |
| `rm` | Remove files/directories | `-rf` (recursive) |
| `touch` | Create empty file | |
| `mv` | Move/rename | |
| `cp` | Copy file | |
| `cat` | Display file contents | `-n` (line numbers), `-b` (non-blank), `-s` (squeeze) |
| `info` | File information | |
| `dir` | List directory | |

### System Monitoring

| Command | Description |
| :--- | :--- |
| `ps` | Process status |
| `sys` | System information |
| `mu` | Memory usage |
| `du` | Disk usage |
| `ut` | Uptime |

### Text Processing

| Command | Description | Flags |
| :--- | :--- | :--- |
| `grep` | Search text using regex | `-f` (ignore case), `-r` (recursive), `-v` (invert) |
| `head` | Display first lines | `-n` (count) |
| `tail` | Display last lines | `-n` (count) |

### Networking

| Command | Description |
| :--- | :--- |
| `ping` | Ping a host |

### Archives

| Command | Description |
| :--- | :--- |
| `zip` | Create ZIP archive |
| `unzip` | Extract ZIP archive |

### Shell Utilities

| Command | Description | Flags |
| :--- | :--- | :--- |
| `echo` | Display text | `-n` (no newline), `-e` (escapes) |
| `boom` | Explosive entrance | |
| `exit` | Exit shell | |

---

## Contributing

Contributions are welcome! Please submit a PR or open an issue.

---

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

---

## Acknowledgments

- [urfave/cli](https://github.com/urfave/cli) - CLI framework
- [gopsutil](https://github.com/shirou/gopsutil) - Cross-platform system monitoring
- [fastping](https://github.com/tatsushid/go-fastping) - Network diagnostics