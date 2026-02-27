module github.com/Masralai/gosh

go 1.25

require (
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/tatsushid/go-fastping v0.0.0-20160109021039-d7bb493dee3e
	github.com/urfave/cli/v3 v3.6.1
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/tklauser/go-sysconf v0.3.16 // indirect
	github.com/tklauser/numcpus v0.11.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
)

// go mod edit -go <version>
// go mod tidy

// during version upgrade
// Clear the build cache, go clean -cache or
