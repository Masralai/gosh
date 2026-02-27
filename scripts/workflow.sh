#!/bin/bash
set -e

echo "--- Starting Pre-Deployment Checks ---"

# 1. Formatting
echo "Checking formatting..."
go fmt ./...

# 2. Tidy Modules
echo "Cleaning up go.mod..."
go mod tidy

# 3. Linting (The part that failed in CI)
echo "Running golangci-lint..."
# If you don't have it: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# Using $(go env GOPATH) ensures it finds the binary regardless of your OS setup
$(go env GOPATH)/bin/golangci-lint run ./...

# 4. Security Scan
echo "Running gosec..."
# If you don't have gosec installed: go install github.com/securego/gosec/v2/cmd/gosec@latest
$(go env GOPATH)/bin/gosec ./...

# 5. Vulnerability Check
echo "Running govulncheck..."
go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# 6. Tests
echo "Running tests..."
go test -v ./...

# 7. Build
echo "Attempting build..."
go build -o gosh-test.exe ./gosh.go

echo "--- All checks passed! Ready to push. ---"
rm gosh-test.exe