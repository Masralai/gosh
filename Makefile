.PHONY: build test lint vet gosec clean

build:
	go build -o gosh .

test:
	go test -race ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

gosec:
	gosec ./...

clean:
	rm -f gosh