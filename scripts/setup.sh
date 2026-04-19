#!/bin/bash

set -e

go env -w GOTOOLCHAIN=local+path
go env -w GO111MODULE=on

GOBIN=/usr/local/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
GOBIN=/usr/local/bin go install github.com/securego/gosec/v2/cmd/gosec@latest
GOBIN=/usr/local/bin go install golang.org/x/vuln/cmd/govulncheck@latest

go mod download
go mod tidy
go mod verify