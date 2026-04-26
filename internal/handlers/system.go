package handlers

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/urfave/cli/v3"
	"os"
	"runtime"
)

func Ps() *cli.Command {
	return &cli.Command{
		Name:  "ps",
		Usage: "process status",
		Action: func(ctx context.Context, c *cli.Command) error {
			ps, err := os.ReadDir("/proc")
			if err != nil {
				return fmt.Errorf("failed to show process status: %v", err)
			}
			fmt.Println(ps)
			return nil
		},
	}
}
func Ut() *cli.Command {
	return &cli.Command{
		Name:      "ut",
		Usage:     "System Uptime",
		UsageText: "cli ut",
		Action: func(ctx context.Context, c *cli.Command) error {
			h, err := host.Uptime()
			if err != nil {
				return fmt.Errorf("failed to fetch system uptime: %v", err)
			}
			fmt.Printf("%f min\n", float64(h)/60)
			return nil

		},
	}
}
func Sys() *cli.Command {
	return &cli.Command{
		Name:      "sys",
		Usage:     "System info",
		UsageText: "cli sys",
		Action: func(ctx context.Context, c *cli.Command) error {

			hs, err := os.Hostname()
			if err != nil {
				return fmt.Errorf("failed to get system hostname: %v", err)
			}
			fmt.Println("hostname:", hs)

			fmt.Println("number of available cpu:", runtime.NumCPU())

			platform, family, version, err := host.PlatformInformation()
			if err != nil {
				return fmt.Errorf("failed to get platform info: %v", err)
			}
			fmt.Println("platform:", platform)
			fmt.Println("family:", family)
			fmt.Println("version:", version)

			kv, err := host.KernelVersion()
			if err != nil {
				return fmt.Errorf("failed to get Kernel Version: %v", err)
			}
			fmt.Println(kv)

			ka, err := host.KernelArch()
			if err != nil {
				return fmt.Errorf("failed to get Kernel Version: %v", err)
			}
			fmt.Println(ka)

			fmt.Println("go version:", runtime.Version())

			ps, err := process.Pids()
			if err != nil {
				return fmt.Errorf("failed to system hostname: %v", err)
			}
			fmt.Println("processes running", ps)

			return nil
		},
	}
}
func Mu() *cli.Command {
	return &cli.Command{
		Name:      "mu",
		Usage:     "Display Free and Used Memory",
		UsageText: "cli mu",
		Action: func(ctx context.Context, c *cli.Command) error {
			v, err := mem.VirtualMemory()
			if err != nil {
				return fmt.Errorf("failed to display memmory usage:%v", err)
			}
			fmt.Printf("Total: %v, Free: %v , UsedPercent: %f%%\n", v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)
			fmt.Println(v.String())

			return nil
		},
	}
}
func Du() *cli.Command {
	return &cli.Command{
		Name:      "du",
		Usage:     "Display disk used",
		UsageText: "cli du <path>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: du <pathname>")
			}

			d, err := disk.Usage(c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to display disk used:%v", err)
			}
			fmt.Printf("Total: %dMB, Free: %dMB, UsedPercent: %f%%\n", d.Total/1024/1024, d.Free/1024/1024, d.UsedPercent)
			fmt.Println("disk usage:", d.String())

			return nil
		},
	}
}
func Kill() *cli.Command {
	return &cli.Command{
		Name:      "kill",
		Usage:     "Terminate Processes using process id",
		UsageText: "cli kill <processname>",
		Action: func(ctx context.Context, c *cli.Command) error {
			ps, err := process.Processes()
			if err != nil {
				return fmt.Errorf("Failed to terminate process:%v", err)
			}
			for _, p := range ps {
				n, err := p.Name()
				if err != nil {
					continue
				}
				if n == c.Args().Get(0) {
					return p.Kill()
				}
			}
			return fmt.Errorf("process not found")
		},
	}
}
