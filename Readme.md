# GoSh (Go-Shell)

[![Go Report Card](https://goreportcard.com/badge/github.com/Masralai/gosh)](https://goreportcard.com/report/github.com/Masralai/gosh)
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
- [Acknowledgments](#acknowledgments)

---

## Key Features

- **Robust Command Engine:** Utilizes `urfave/cli/v3` for sophisticated argument parsing and flag management.
- **Enhanced Security:** Path traversal guards, directory traversal checks, and confirmation prompts for destructive actions.
- **Portability:** Works natively or containerized with Docker.
- **Interactive Experience:** Readline-based REPL with quote-aware argument parsing for builtins, persistent command history, tab completion for builtin command names/aliases, and host-shell fallback for unknown commands.
- **System Monitoring:** Built-in tools for monitoring CPU, memory, disk usage, and process management.
- **Network Utilities:** Integrated network diagnostic tools including ping functionality (ICMP ping may require elevated privileges on some systems).
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

| Component            | Description                                      |
| -------------------- | ------------------------------------------------ |
| `cmd/gosh/main.go`   | Entry point, sets up urfave/cli and command loop |
| `internal/handlers/` | All command implementations                      |

### Command Dispatcher

The core engine is built on `urfave/cli/v3`. It maps user input to specific handler functions in `internal/handlers/`. New commands can be added by implementing a handler function and registering it in `commands.go`.

---

## Prerequisites

- **Go:** Version 1.25 or higher
- **Docker:** Optional, for containerized execution

---

## Installation

### Install script (Linux, macOS, Windows via Git Bash)

```bash
curl -fsSL https://raw.githubusercontent.com/Masralai/gosh/main/install.sh | sh
```

Installs to `~/.local/bin` after checking the release checksum. Set `GOSH_BIN` to change the directory.

### Native (Recommended)

1. Clone the repository:

```bash
git clone https://github.com/Masralai/gosh
cd gosh
```

1. Build:

```bash
make build
# or
go build -o gosh ./cmd/gosh
```

1. Run:

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

### Shell Utilities

| Command   | Description                                      | Flags / Notes |
| :-------- | :----------------------------------------------- | :------------ |
| `help`    | Show help for built-in commands                  | `help [command]` |
| `history` | Display or clear persisted command history       | `-c` (clear history) |
| `env`     | Display environment variables                    | `env [name...]` |
| `which`   | Locate a command (GoSh builtin/alias or host PATH) | `which <name> [name...]` |
| `echo`    | Display a line of text                           | `-n` (no newline), `-e` (escapes); prints all arguments |
| `clear`   | Clear the terminal screen                        |               |
| `cd`      | Change directory                                 | `cd <path>`   |
| `pwd`     | Print working directory                          |               |
| `boom`    | Explosive entrance                               |               |
| `exit`    | Exit shell                                       |               |

### File Operations

| Command | Description                                              | Flags / Notes |
| :------ | :------------------------------------------------------- | :------------ |
| `ls`    | List directory contents                                  | `-a` (hidden), `-l` (long), `-h` (human sizes with `-l`), `-R` (recursive), `-S` (sort by size), `-t` (sort by mtime) |
| `mkdir` | Create directories                                       | `-p` (create parents; no error if exists) |
| `rm`    | Remove files or directories                              | `-r`/`-R` (recursive), `-f` (force, ignore missing, no prompt). `rm -rf` is Unix-like (no confirmation prompt). |
| `touch` | Create file or update its timestamp (no truncate)        | `touch <file> [file...]` |
| `mv`    | Move/rename file or directory                            | `mv <source> <destination>` |
| `cp`    | Copy files and directories                               | `-r`/`-R` (recursive directory copy) |
| `dir`   | Display directory contents (human-readable)              | `dir [path]` |
| `cat`   | Display file contents (supports multiple files)          | `-n` (number all lines), `-b` (number non-blank, overrides `-n`), `-s` (squeeze blank lines) |
| `info`  | Display file information                                 | `info <path>` |
| `find`  | Search for files by name pattern (basic recursive search) | `find <path> -name <pattern>` |

### System Monitoring

| Command | Description                        |
| :------ | :--------------------------------- |
| `ps`    | Display process status             |
| `sys`   | Display system information         |
| `mu`    | Display memory usage               |
| `du`    | Display disk usage for a path      |
| `ut`    | Display system uptime              |
| `kill`  | Terminate a process by PID         |

### Text Processing

| Command | Description                                | Flags / Notes |
| :------ | :----------------------------------------- | :------------ |
| `grep`  | Search for a pattern in files              | `-i` (ignore case, primary), `-n` (line numbers), `-r` (recursive), `-v` (invert). `-f` is kept as a compatibility alias for ignore-case (not pattern-file). |
| `head`  | Display the beginning of a file            | `-n <lines>` (default 10) |
| `tail`  | Display the end of a file                  | `-n <lines>` (default 10) |
| `wc`    | Count lines, words, and bytes              | `-l`, `-w`, `-c` (default: all three) |
| `sort`  | Sort lines of text                         | `-r` (reverse) |

### Networking

| Command | Description |
| :------ | :---------- |
| `ping`  | Send ICMP echo requests to a host (may require elevated privileges on some systems; GoSh reports clear errors when unavailable) |

### Archives

| Command | Description         | Usage / Notes                                      |
| :------ | :------------------ | :------------------------------------------------- |
| `zip`   | Create a ZIP archive from files | `zip <archive> <file1> [file2] ...` (adds `.zip` if omitted) |
| `unzip` | Extract files from a ZIP archive | `unzip <archive.zip> <destination>`               |

---

## Notes on Behavior Changes

**Interactive/QOL milestone**
- Interactive builtin input now uses quote-aware argument parsing (supports quoted paths/arguments with spaces).
- Command history is persisted across sessions (configurable via `GOSH_HISTORY`; defaults to a per-user config/history location).
- Tab completion suggests builtin command names and aliases (not file/path arguments).
- `ls` implements `-R` (recursive) and `-S` (sort by size) as advertised.
- `cat -b` numbers non-blank lines (and overrides `-n` when both are provided).
- `echo` prints all provided arguments.
- `dir` prints a human-readable listing (not internal structures) and accepts an optional path.
- `kill` terminates a process by PID (`kill <pid>`).
- `touch` updates timestamps without truncating existing files.
- `zip` avoids adding a duplicate `.zip` extension when one is already provided.
- `ping` keeps ICMP-based behavior but reports clearer cross-platform/privilege-related errors (ICMP may require elevation on some systems).

**Broadening milestone**
- Added builtins: `help`, `history` (with `-c`), `env`, `which`, `find` (basic `-name`), `wc`, and `sort`.
- `ls` adds `-l` (long listing), `-h` (human-readable sizes with `-l`), and `-t` (sort by modification time).
- `mkdir` adds `-p` to create parent directories and succeed if the directory already exists.
- `cp` adds `-r`/`-R` for recursive directory copies.
- `rm` now uses familiar Unix-like semantics: `-r`/`-R` removes directories recursively and `-f` forces removal (ignores missing targets, no prompts). `rm -rf` no longer prompts for confirmation (behavior change from the previous milestone).
- `grep` adds `-i` (primary ignore-case) and `-n` (line numbers). `-f` is kept as a compatibility alias for ignore-case in this milestone (not pattern-file behavior).
- `cat` accepts multiple files and applies `-n`/`-b`/`-s` consistently.
- `which` reports GoSh builtins/aliases and best-effort host PATH resolution.

---

## Contributing

Contributions are welcome! Please submit a PR or open an issue.

---

## Acknowledgments

- [urfave/cli](https://github.com/urfave/cli) - CLI framework
- [gopsutil](https://github.com/shirou/gopsutil) - Cross-platform system monitoring
- [fastping](https://github.com/tatsushid/go-fastping) - Network diagnostics
- [readline](https://github.com/chzyer/readline) - Interactive line editing, history, and completion support
