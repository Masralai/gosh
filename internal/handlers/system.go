package handlers

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/urfave/cli/v3"
)

func Ps() *cli.Command {
	return &cli.Command{
		Name:      "ps",
		Usage:     "Display process status",
		UsageText: "ps",
		Action: func(ctx context.Context, c *cli.Command) error {
			ps, err := process.Processes()
			if err != nil {
				return fmt.Errorf("failed to show process status: %v", err)
			}
			for _, p := range ps {
				name, err := p.Name()
				if err != nil {
					continue
				}
				fmt.Printf("%6d  %s\n", p.Pid, name)
			}
			return nil
		},
	}
}
func Ut() *cli.Command {
	return &cli.Command{
		Name:      "ut",
		Usage:     "Display system uptime",
		UsageText: "ut",
		Action: func(ctx context.Context, c *cli.Command) error {
			h, err := host.Uptime()
			if err != nil {
				return fmt.Errorf("failed to fetch system uptime: %v", err)
			}
			fmt.Printf("%.1f min\n", float64(h)/60)
			return nil
		},
	}
}
func Sys() *cli.Command {
	return &cli.Command{
		Name:      "sys",
		Usage:     "Display system information",
		UsageText: "sys",
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
				return fmt.Errorf("failed to get kernel version: %v", err)
			}
			fmt.Println("kernel:", kv)

			ka, err := host.KernelArch()
			if err != nil {
				return fmt.Errorf("failed to get kernel architecture: %v", err)
			}
			fmt.Println("arch:", ka)

			fmt.Println("go version:", runtime.Version())

			ps, err := process.Pids()
			if err != nil {
				return fmt.Errorf("failed to list processes: %v", err)
			}
			fmt.Println("processes:", len(ps))

			return nil
		},
	}
}
func Mu() *cli.Command {
	return &cli.Command{
		Name:      "mu",
		Usage:     "Display memory usage",
		UsageText: "mu",
		Action: func(ctx context.Context, c *cli.Command) error {
			v, err := mem.VirtualMemory()
			if err != nil {
				return fmt.Errorf("failed to display memory usage: %v", err)
			}
			fmt.Printf("Memory  %s total  %s free  %.0f%% used  %s\n", humanBytes(v.Total), humanBytes(v.Free), v.UsedPercent, meter(v.UsedPercent))
			return nil
		},
	}
}
func Du() *cli.Command {
	return &cli.Command{
		Name:      "du",
		Usage:     "Display disk usage for a path",
		UsageText: "du <path>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: du <path>")
			}

			d, err := disk.Usage(c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to display disk usage: %v", err)
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
		Usage:     "Terminate a process by PID",
		UsageText: "kill <pid>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: kill <pid>")
			}
			pidArg := c.Args().Get(0)
			pid, err := strconv.Atoi(pidArg)
			if err != nil || pid <= 0 {
				return fmt.Errorf("invalid pid %q: must be a positive integer", pidArg)
			}
			p, err := process.NewProcess(int32(pid))
			if err != nil {
				return fmt.Errorf("process %d not found: %v", pid, err)
			}
			if err := p.Kill(); err != nil {
				return fmt.Errorf("failed to terminate process %d: %v", pid, err)
			}
			return nil
		},
	}
}

func humanBytes(n uint64) string {
	const gb = 1024 * 1024 * 1024
	if n >= gb {
		return fmt.Sprintf("%.1f GB", float64(n)/gb)
	}
	return fmt.Sprintf("%.0f MB", float64(n)/1024/1024)
}

func meter(percent float64) string {
	filled := int(percent/10 + 0.5)
	if filled < 0 {
		filled = 0
	}
	if filled > 10 {
		filled = 10
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", 10-filled)
}
