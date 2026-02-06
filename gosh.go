package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	// "os/exec"
	// "path/filepath"
	//"net/http"

	// "github.com/shirou/gopsutil"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/urfave/cli/v3"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\033[H\033[2J") //clears terminal
	fmt.Println(`
       ^  ^  ^   ^      ___I_      ^  ^   ^  ^  ^   ^  ^
      /|\/|\/|\ /|\    /\-_--\    /|\/|\ /|\/|\/|\ /|\/|\
      /|\/|\/|\ /|\   /  \_-__\   /|\/|\ /|\/|\/|\ /|\/|\
      /|\/|\/|\ /|\   |[]| [] |   /|\/|\ /|\/|\/|\ /|\/|\
	
	Welcome to GoSh
	To get started: cli -h or cli {command} -h
	  `)
	for scanner.Scan() {

		root := &cli.Command{
			Name:    "GoSh",
			Version: "v1.2.0",
			Authors: []any{
				"Name:  Devdeep Paul",
				"Email: devv.v4828@gmail.com",
			},
			Usage:     "A cozy, cross-platform shell environment built with Go.",
			UsageText: "cli [command] [arguments]",

			//basic commands
			Commands: []*cli.Command{
				{
					Name:      "boom",
					Usage:     "make an explosive entrance",
					UsageText: "cli boom",
					Action: func(context.Context, *cli.Command) error {
						fmt.Println("boom! I say!")
						return nil
					},
				},
				{
					Name:      "echo",
					Usage:     "Display Text",
					UsageText: "cli echo <text>",
					Action: func(ctx context.Context, c *cli.Command) error {

						fmt.Println(c.Args().Get(0))

						return nil
					},
				},
				{
					Name:      "cd",
					Usage:     "Change Directory",
					UsageText: "cli cd <path>",

					Action: func(ctx context.Context, c *cli.Command) error {
						os.Chdir(c.Args().Get(0))
						fmt.Println(os.Getwd())
						return nil
					},
				},
				{
					Name:      "pwd",
					Usage:     "Print Working Directory",
					UsageText: "cli pwd",
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println(os.Getwd())
						return nil
					},
				},
				{
					Name:      "exit",
					Usage:     "exit",
					UsageText: "cli exit",
					Action: func(context.Context, *cli.Command) error {
						os.Exit(0)
						return nil
					},
				},
				{
					Name:      "ls",
					Usage:     "List Directory Contents",
					UsageText: "cli ls",
					Action: func(context.Context, *cli.Command) error {
						fmt.Println(os.ReadDir("./"))
						return nil
					},
				},
				{
					Name:      "mkdir",
					Usage:     "Make Directories",
					UsageText: "cli mkdir <path>",
					Action: func(ctx context.Context, c *cli.Command) error {
						os.Mkdir(c.Args().Get(0), 64)
						return nil
					},
				},
				{
					Name:      "rm",
					Usage:     "Remove Files or Directories",
					UsageText: "cli rm <filename> or cli rm -rf <path>",

					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "rf",
							Aliases: []string{"rf", "r"},
							Usage:   "recursive delete",
						},
					},
					Action: func(ctx context.Context, c *cli.Command) error {
						if c.Args().Len() == 0 {
							// return cli.Exit("no file specified",14)
							fmt.Println("No file specified")
							return nil
						}
						if c.Bool("rf") {
							fmt.Println("are you sure you want to trigger recursive deletion ? y/n")

							response := ""
							fmt.Scanln(&response)
							if response == "y" || response == "Y" {
								os.RemoveAll(c.Args().Get(0))
							} else {
								fmt.Println("Aborted")
							}

						} else {
							// for i:=0 ;i<len(c.Args().Slice());i++ {
							// 		os.Remove(c.Args().Get(i))
							// }
							for _, filename := range c.Args().Slice() {
								err := os.Remove(filename)
								if err != nil {
									fmt.Printf("Error deleting %s: %v\n", filename, err)
								}
							}
						}
						return nil
					},
				},

				{
					Name:      "touch",
					Usage:     "change file timestamps or create an empty file",
					UsageText: "cli touch <filename>",

					Action: func(ctx context.Context, c *cli.Command) error {
						os.Create(c.Args().Get(0))
						return nil
					},
				},
				{
					Name:      "mv",
					Usage:     "move/rename file",
					UsageText: "cli rn <filename> <filename>",

					Action: func(ctx context.Context, c *cli.Command) error {
						os.Rename(c.Args().Get(0), c.Args().Get(1))
						return nil
					},
				},
				{
					Name:      "dir",
					Usage:     "directory contents",
					UsageText: "cli dir <path>",

					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println(os.ReadDir(c.Args().Get(0)))
						return nil
					},
				},
				{
					Name:      "cat",
					Usage:     "read contents",
					UsageText: "cli cat <filename>",
					Action: func(ctx context.Context, c *cli.Command) error {
						data, err := os.ReadFile(c.Args().Get(0))
						if err != nil {
							log.Fatal(err)
						}
						os.Stdout.Write(data)
						return nil
					},
				},
				{
					Name:      "info",
					Usage:     "file info",
					UsageText: "cli info <filename>",
					Action: func(ctx context.Context, c *cli.Command) error {
						s, err := os.Stat(c.Args().Get(0))
						if err != nil {
							return err
						} else {
							fmt.Printf("File: %s\n", s.Name())
							fmt.Printf("Size: %d bytes\n", s.Size())
							// On Windows, Go’s os.FileMode doesn't strictly report the "Execute" bit because Windows determines
							// what is executable based on the file extension (.exe, .bat), not a permission bit like Linux does.
							// In Go's eyes on Windows, most files show up as rw.
							fmt.Printf("Read(r),Write(w),Execute(x) | Mode: %s\n", s.Mode())
							fmt.Printf("Last Modified: %s\n", s.ModTime().Format("2006-01-02 15:04:05"))
							fmt.Printf("Directory?: %v\n", s.IsDir())
						}

						return nil
					},
				},

				//////////////////////////////////////////////////////////////////////////////
				//system monitoring
				{
					Name:  "ps",
					Usage: "process status",
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println(os.ReadDir("/proc"))
						return nil
					},
				},
				{
					Name:      "ut",
					Usage:     "System Uptime",
					UsageText: "cli ut",
					Action: func(ctx context.Context, c *cli.Command) error {
						h, _ := host.Uptime()
						fmt.Printf("%f min\n", float64(h/60))
						return nil
					},
				},
				{
					Name:      "sys",
					Usage:     "System info",
					UsageText: "cli sys",
					Action: func(ctx context.Context, c *cli.Command) error {

						o, _ := os.Hostname()
						fmt.Println("hostname:", o)
						fmt.Println("number of available cpu:", runtime.NumCPU())
						fmt.Println(host.PlatformInformation())
						fmt.Println(host.KernelVersion())
						fmt.Println(host.KernelArch())
						fmt.Println("go version:", runtime.Version())
						p, _ := process.Pids()
						fmt.Println("processes running", p)

						return nil
					},
				},
				{
					Name:      "mu",
					Usage:     "Display Free and Used Memory",
					UsageText: "cli mu",
					Action: func(ctx context.Context, c *cli.Command) error {
						v, _ := mem.VirtualMemory()
						fmt.Printf("Total: %v, Free: %v , UsedPercent: %f%%\n", v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)
						fmt.Println(v.String())

						return nil
					},
				},
				{
					Name:      "du",
					Usage:     "Display disk used",
					UsageText: "cli du <path>",
					Action: func(ctx context.Context, c *cli.Command) error {
						if c.Args().Len() == 0 {
							println("provide filesystem path such as /")
						} else {

							d, _ := disk.Usage(c.Args().Get(0))
							fmt.Printf("Total: %dMB, Free: %dMB, UsedPercent: %f%%\n", d.Total/1024/1024, d.Free/1024/1024, d.UsedPercent)
							fmt.Println("disk usage:", d.String())
						}
						return nil
					},
				},
				{
					Name:      "kill",
					Usage:     "Terminate Processes using process id",
					UsageText: "cli kill <processname>",
					Action: func(ctx context.Context, c *cli.Command) error {
						ps, err := process.Processes()
						if err != nil {
							return err
						}
						for _, p := range ps {
							n, err := p.Name()
							if err != nil {
								return err
							}
							if n == c.Args().Get(0) {
								return p.Kill()
							}
						}
						return fmt.Errorf("process not found")
					},
				},
				//////////////////////////////////////////////////////////////////////////////
				//Text Processing
				{
					Name:      "grep",
					Usage:     "Search Text Using Patterns",
					UsageText: "cli grep 'pattern' <filename>",
					Action: func(ctx context.Context, c *cli.Command) error {
						return nil
					},
				},
				{
					Name:      "cli head",
					Usage:     "Display the beginning of a file",
					UsageText: "cli head <filename>",
					Action: func(ctx context.Context, c *cli.Command) error {
						return nil
					},
				},
				{
					Name:      "tail",
					Usage:     "Display Last Part of Files",
					UsageText: "cli tail <filename>",
					Action: func(ctx context.Context, c *cli.Command) error {
						return nil
					},
				},

				//////////////////////////////////////////////////////////////////////////////
				//Networking
				{
					Name:      "ping",
					Usage:     "Send Request to Network Hosts",
					UsageText: "cli ping <hostname>",
					Action: func(ctx context.Context, c *cli.Command) error {
						return nil
					},
				},
				{
					Name:      "curl",
					Usage:     "Transfer a URL",
					UsageText: "cli curl http://example.com/file.txt",
					Action: func(ctx context.Context, c *cli.Command) error {
						return nil
					},
				},
			},

			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "ginger-crouton",
					Usage: "is it in the soup?",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				if !cmd.Bool("ginger-crouton") {
					fmt.Println("invalid command")
				}
				return nil
			},
		}
		if err := root.Run(context.Background(), strings.Fields(scanner.Text())); err != nil { //ignores excess spaces and tabs
			log.Fatal(err)
		}

	}

}
