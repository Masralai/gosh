#!/bin/bash
set -e

echo "--- Starting Pre-Deployment Checks ---"

# 1. Formatting
echo "Checking formatting..."
go fmt ./...

# 2. Tidy Modules
echo "Cleaning up go.mod..."
go mod tidy

# 3. Linting
echo "Running golangci-lint..."
# Check common paths for golangci-lint
if command -v golangci-lint &> /dev/null; then
    golangci-lint run ./...
elif [ -x "$(go env GOPATH)/bin/golangci-lint" ]; then
    $(go env GOPATH)/bin/golangci-lint run ./...
elif [ -x "/usr/local/bin/golangci-lint" ]; then
    /usr/local/bin/golangci-lint run ./...
else
    echo "Warning: golangci-lint not found, skipping..."
fi

# 4. Security Scan
echo "Running gosec..."
# Check common paths for gosec
if command -v gosec &> /dev/null; then
    gosec ./...
elif [ -x "$(go env GOPATH)/bin/gosec" ]; then
    $(go env GOPATH)/bin/gosec ./...
elif [ -x "/usr/local/bin/gosec" ]; then
    /usr/local/bin/gosec ./...
else
    echo "Warning: gosec not found, skipping..."
fi

# 5. Vulnerability Check
echo "Running govulncheck..."
go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# 6. Tests
echo "Running tests..."
go test -v ./...

# 7. Build
echo "Attempting build..."
go build -o gosh ./cmd/gosh

echo "--- All checks passed! Ready to push. ---"
rm -f gosh
