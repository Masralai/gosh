# GoSh (Go-Shell)

[![Go Report Card](https://goreportcard.com/badge/github.com/Masralai/gosh)](https://goreportcard.com/report/github.com/Masralai/gosh)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Masralai/gosh)](https://github.com/Masralai/gosh)

GoSh is a modular, interactive command-line interface built in Go. It provides a containerized environment for secure command execution, system monitoring, and file management.

---

## Table of Contents

- [Key Features](#key_features)
- [Architecture Overview](#architecture-overview)
- [Prerequisites](#prerequisites)
- [Installation and Execution](#installation-and-execution)
- [Command Reference](#command-reference)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

---

## Key Features

- **Robust Command Engine:** Utilizes `urfave/cli/v3` for sophisticated argument parsing and flag management.
- **Enhanced Security:** Custom file operation logic includes recursive guards and mandatory confirmation prompts for destructive actions.
- **Portability:** Fully containerized architecture using Docker ensures a consistent, zero-install execution environment.
- **Interactive Experience:** Integrated command scanning provides a fluid terminal interaction with persistent history.
- **System Monitoring:** Built-in tools for monitoring CPU, memory, disk usage, and process management.
- **Network Utilities:** Integrated network diagnostic tools including ping functionality.
- **Archive Management:** Support for creating and extracting ZIP archives with security checks against decompression bombs.

---

## Architecture Overview

GoSh is designed with a modular architecture that separates the command interface from the underlying system logic.

### Command Dispatcher

The core engine is built on `urfave/cli/v3`. It acts as a dispatcher that maps user input to specific Go functions. This allows for easy extensibility—new commands can be added by simply registering them in the `root` command structure in `gosh.go`.

### Containerization Strategy

The project leverages Docker to provide an isolated execution environment. This ensures that:

1. **Security:** Commands executed within the shell are isolated from the host system.
2. **Environment Consistency:** Dependencies and OS-level configurations are standardized across all deployments.

---

## Prerequisites

- **Go:** Version 1.25 or higher (for native builds).
- **Docker:** Docker Desktop or Docker Engine (recommended for isolated execution).

---

## Installation and Execution

### Docker Execution (Recommended)

Running GoSh via Docker is the most efficient method as it avoids modifying your local system environment.

```powershell
docker compose run gosh
```

This command automatically builds the necessary images and initiates an interactive TTY session.

### Native Local Build

To run GoSh directly on your host machine, follow these steps:

1. **Initialize Module:**

   ```bash
   go mod init go-cli
   ```

2. **Install Dependencies:**

   ```bash
   go get github.com/urfave/cli/v3
   ```

3. **Launch Application:**

   ```bash
   go run gosh.go
   ```

---

## Command Reference

| Command | Category | Description | Primary Flags |
| :--- | :--- | :--- | :--- |
| `ls` | File Ops | List directory contents | `-a`, `-R`, `-S` |
| `cd` | File Ops | Change working directory | |
| `rm` | File Ops | Remove files or directories | `-rf` (Recursive) |
| `sys` | System | Display system information | |
| `mu` | System | Memory usage stats | |
| `du` | System | Disk usage stats | |
| `grep` | Text | Search text using regex | `-f`, `-r`, `-v` |
| `ping` | Network | Network host diagnostic | |
| `zip` | Storage | Create ZIP archive | |
| `unzip`| Storage | Extract ZIP archive | |

---

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

---

## Acknowledgments

- [urfave/cli](https://github.com/urfave/cli) - For the robust CLI framework.
- [gopsutil](https://github.com/shirou/gopsutil) - For cross-platform system monitoring.
- [fastping](https://github.com/tatsushid/go-fastping) - For network diagnostic capabilities.
