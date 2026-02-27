package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	// "log"
	"archive/zip"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	// "net/http"
	"net"
	"time"
	// "github.com/shirou/gopsutil"
	// "github.com/tklauser/go-sysconf"
	//"github.com/joho/godotenv"
	// curl "github.com/andelf/go-curl"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/tatsushid/go-fastping"
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
			},
			Usage:     "A cozy, cross-platform shell environment built with Go.",
			UsageText: "cli [command] [arguments]",

			//basic commands
			Commands: []*cli.Command{
				{
					Name:            "cli",
					SkipFlagParsing: true,
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println("cli", c.Args())
						return nil
					},
				},
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
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "n",
							Usage: "Don't add a new line at the end",
						},
						&cli.StringFlag{
							Name:  "e",
							Usage: "Allow special characters like \n for new lines",
						},
					},
					Action: func(ctx context.Context, c *cli.Command) error {
						if c.Args().Len() == 0 {
							return fmt.Errorf("usage: echo <text>")
						}
						fmt.Println(c.Args().Get(0))
						return nil
					},
				},
				{
					Name:      "cd",
					Usage:     "Change Directory",
					UsageText: "cli cd <path>",

					Action: func(ctx context.Context, c *cli.Command) error {
						if err := os.Chdir(c.Args().Get(0)); err != nil {
							return fmt.Errorf("failed to changed dir: %v", err)
						}

						wd, err := os.Getwd()
						if err != nil {
							return fmt.Errorf("failed to change Directory : %v", err)
						}
						fmt.Println(wd)
						return nil
					},
				},
				{
					Name:      "pwd",
					Usage:     "Print Working Directory",
					UsageText: "cli pwd",
					Action: func(ctx context.Context, c *cli.Command) error {
						wd, err := os.Getwd()
						if err != nil {
							return fmt.Errorf("failed to Print Working Directory : %v", err)
						}
						fmt.Println(wd)

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
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "R",
							Usage: "List subdirectories recursively",
						},
						&cli.StringFlag{
							Name:  "S",
							Usage: "Sort by file size",
						},
						&cli.StringFlag{
							Name:  "a",
							Usage: "Include hidden files",
						},
					},
					Action: func(context.Context, *cli.Command) error {
						rd, err := os.ReadDir("./")
						if err != nil {
							return fmt.Errorf("failed to List Directory Contents: %v", err)
						}
						fmt.Println(rd)

						return nil
					},
				},
				{
					Name:      "mkdir",
					Usage:     "Make Directories",
					UsageText: "cli mkdir <path>",
					Action: func(ctx context.Context, c *cli.Command) error {
						if c.Args().Len() == 0 {
							return fmt.Errorf("usage: mkdir <pathname>")
						}
						err := os.Mkdir(c.Args().Get(0), 64)
						if err != nil {
							return fmt.Errorf("failed to created directory: %v", err)
						}
						fmt.Printf("directory created %s", c.Args().Get(0))
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
							return fmt.Errorf("usage: rm <filename> || usage: rm -rf <filename>")

						}
						if c.Bool("rf") {
							fmt.Println("are you sure you want to trigger recursive deletion ? y/n")

							response := ""
							if _, err := fmt.Scanln(&response); err != nil {
								return fmt.Errorf("failed to get user input:%v", err)
							}
							if response == "y" || response == "Y" {
								if err := os.RemoveAll(c.Args().Get(0)); err != nil {
									return fmt.Errorf("Failed recursive delete:%v", err)
								}
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
									return fmt.Errorf("Error deleting %s: %v\n", filename, err)
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
						if c.Args().Len() == 0 {
							return fmt.Errorf("usage: touch <filename>")
						}
						_, err := os.Create(c.Args().Get(0))
						if err != nil {
							return fmt.Errorf("failed to create file :%v", err)
						}
						fmt.Printf("file created : %s\n", c.Args().Get(0))

						return nil
					},
				},
				{
					Name:      "mv",
					Usage:     "move/rename file",
					UsageText: "cli mv <src> <dest>",

					Action: func(ctx context.Context, c *cli.Command) error {
						if c.Args().Len() < 2 {
							return fmt.Errorf("usage: mv <source> <destination>")
						}
						src := c.Args().Get(0)
						dest := c.Args().Get(1)
						err := os.Rename(src, dest)
						if err != nil {
							return fmt.Errorf("failed to move %q :%v", src, err)
						}

						fmt.Printf("moved %q to %q\n", src, dest)
						return nil
					},
				},
				{
					Name:      "cp",
					Usage:     "Copy Files and Directories within",
					UsageText: "cli cp <src> <dest>",

					Action: func(ctx context.Context, c *cli.Command) error {
						if c.Args().Len() < 2 {
							return fmt.Errorf("usage: cp <source> <destination>")
						}
						src := c.Args().Get(0)
						dest := c.Args().Get(1)
						srcDir := os.DirFS(src)
						err := os.CopyFS(dest, srcDir)
						if err != nil {
							return fmt.Errorf("failed to copy %q :%v", src, err)
						}
						fmt.Printf("copied %q to %q\n", src, dest)
						return nil
					},
				},
				{
					Name:      "dir",
					Usage:     "directory contents",
					UsageText: "cli dir <path>",

					Action: func(ctx context.Context, c *cli.Command) error {
						if c.Args().Len() == 0 {
							return fmt.Errorf("usage: dir <pathname>")
						}
						dir, err := os.ReadDir(c.Args().Get(0))
						if err != nil {
							return fmt.Errorf("failed to display dir contents: %v", err)
						}
						fmt.Println(dir)
						return nil
					},
				},
				{
					Name:      "cat",
					Usage:     "read contents",
					UsageText: "cli cat <filename>",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "n",
							Usage: "Add numbers to each line",
						},
						&cli.StringFlag{
							Name:  "b",
							Usage: "Add numbers only to lines with text",
						},
						&cli.StringFlag{
							Name:  "s",
							Usage: "Remove extra empty lines",
						},
					},
					Action: func(ctx context.Context, c *cli.Command) error {
						data, err := os.ReadFile(c.Args().Get(0))
						if err != nil {
							// log.Fatal(err)
							return fmt.Errorf("failed to read file contents: %v", err)
						}
						if _, err := os.Stdout.Write(data); err != nil {
							return fmt.Errorf("failed to write data:%v", err)
						}
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
							return fmt.Errorf("failed to get file info:%v", err)
						}
						fmt.Printf("File: %s\n", s.Name())
						fmt.Printf("Size: %d bytes\n", s.Size())
						// On Windows, Go’s os.FileMode doesn't strictly report the "Execute" bit because Windows determines
						// what is executable based on the file extension (.exe, .bat), not a permission bit like Linux does.
						// In Go's eyes on Windows, most files show up as rw.
						fmt.Printf("Read(r),Write(w),Execute(x) | Mode: %s\n", s.Mode())
						fmt.Printf("Last Modified: %s\n", s.ModTime().Format("2006-01-02 15:04:05"))
						fmt.Printf("Directory?: %v\n", s.IsDir())

						return nil
					},
				},

				//////////////////////////////////////////////////////////////////////////////
				//system monitoring
				{
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
				},
				{
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
				},
				{
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
						fmt.Println(host.PlatformInformation())

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
				},
				{
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
				},
				{
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
				},
				{
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
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "f",
							Usage: "Search ignoring case differences (uppercase or lowercase)",
						},
						&cli.BoolFlag{
							Name:  "r",
							Usage: "recursive grep",
						},
						&cli.BoolFlag{
							Name:  "v",
							Usage: "Find lines that do not match the pattern",
						},
					},
					Action: func(ctx context.Context, c *cli.Command) error {
						if c.Args().Len() < 2 {
							return fmt.Errorf("not enough arguments")
						}

						filename := c.Args().Get(1)
						pattern := c.Args().Get(0)
						regObj, err := regexp.Compile(pattern)
						if err != nil {
							return fmt.Errorf("Failed to create regression Object: %v", err)
						}
						//allow access only within the current directory
						root, _ := os.OpenRoot(".")
						defer root.Close()
						file, err := root.Open(filename)
						if err != nil {
							return fmt.Errorf("Failed to open file %v", err)
						}
						defer file.Close()
						scanner := bufio.NewScanner(file)
						if c.Bool("f") {
							return nil
						}
						if c.Bool("r") {

							for scanner.Scan() {
								line := scanner.Text()
								if regObj.MatchString(line) {
									fmt.Printf("%s\n", line)
								}
							}
						}

						if c.Bool("v") {
							for scanner.Scan() {
								line := scanner.Text()
								if !regObj.MatchString(line) {
									fmt.Printf("%s\n", line)
								}
							}
						}
						for scanner.Scan() { //defer before this func since it has hidden scanner.Err()
							line := scanner.Text() //strips the new line chars from the txt file
							if regObj.MatchString(line) {
								fmt.Printf("%s\n", line)
							}
						}

						return nil
					},
				},
				{
					Name:      "head",
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
						p := fastping.NewPinger()
						ra, err := net.ResolveIPAddr("ip4:icmp", c.Args().Get(0))
						if err != nil {
							return fmt.Errorf("Failed to resolve IP Address: %v", err)

						}
						p.AddIPAddr(ra)
						p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
							fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
						}
						p.OnIdle = func() {
							fmt.Println("finish")
						}
						err = p.Run()
						if err != nil {
							return fmt.Errorf("Failed to ping :%v", err)
						}
						return nil
					},
				},
				// {
				// 	Name:      "curl",
				// 	Usage:     "Transfer a URL",
				// 	UsageText: "cli curl http://example.com/file.txt",
				// 	Action: func(ctx context.Context, c *cli.Command) error {
				// 		easy := curl.EasyInit()
				// 		defer easy.Cleanup()
				// 		if easy != nil {
				// 			easy.Setopt(curl.OPT_URL, c.Args().Get(0))
				// 			easy.Perform()
				// 		}
				// 		return nil
				// 	},
				// },

				//////////////////////////////////////////////////////////////////////////////
				//File Compression
				{
					Name:      "zip",
					Usage:     "compress file",
					UsageText: "cli zip <filename>",
					Action: func(ctx context.Context, c *cli.Command) error {
						archive, err := os.Create(c.Args().Get(0) + ".zip")
						if err != nil {
							return fmt.Errorf("failed to create zip archive :%v", err)
						}
						defer archive.Close()

						zipWriter := zip.NewWriter(archive)
						fmt.Println("opening first file...")
						f1, err := os.Open(c.Args().Get(1))
						if err != nil {
							return fmt.Errorf("file error:%v", err)
						}
						defer f1.Close()
						fmt.Println("adding file to the archive...")

						//compression path
						path := c.Args().Get(3)
						w1, err := zipWriter.Create(path)
						if err != nil {
							return fmt.Errorf("Failed to add file to archive:%v", err)
						}

						//copy uncompressed file to achive
						if _, err := io.Copy(w1, f1); err != nil {
							return fmt.Errorf("Failed to copy uncompressed file to achive:%v", err)
						}

						fmt.Println("opening second file ...")
						f2, err := os.Open(c.Args().Get(2))
						if err != nil {
							return fmt.Errorf("failed to open second file:%v", err)
						}
						defer f2.Close()

						//create entry in the zip archive
						w2, err := zipWriter.Create(path)
						if err != nil {
							return fmt.Errorf("Failed to add file to archive:%v", err)
						}

						//copy uncompressed file to achive
						if _, err := io.Copy(w2, f2); err != nil {
							return fmt.Errorf("Failed to copy uncompressed file to achive:%v", err)
						}

						if err := zipWriter.Close(); err != nil {
							return fmt.Errorf("failed to close zipwriter:%v", err)
						}
						return nil
					},
				},
				{
					Name:      "unzip",
					Usage:     "Extract from ZIP archive",
					UsageText: "cli unzip <filename>.zip dest",
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println("opening zip archive")
						filename := c.Args().Get(0)
						archive, err := zip.OpenReader(filename)
						if err != nil {
							return fmt.Errorf("failed to read archive: %v", err)
						}
						defer archive.Close()

						dest := c.Args().Get(1)
						root, err := os.OpenRoot(dest)
						if err != nil {
							return fmt.Errorf("failed to open dest root:%v", err)
						}
						defer root.Close()

						for _, f := range archive.File {
							// filePath := filepath.Join(dest, f.Name)
							filePath:= f.Name

							destAbs, _ := filepath.Abs(dest)
							fileAbs, _ := filepath.Abs(filePath)
							if !strings.HasPrefix(fileAbs, destAbs) {
								return fmt.Errorf("illegal file path:%s", f.Name)
							}

							fmt.Println("unzipping file...", filePath)
							//Decompression Bomb prevention
							if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
								return fmt.Errorf("invalid file path %s", filePath)
							}

							//empty dir
							if f.FileInfo().IsDir() {
								fmt.Println("creating directory")
								//os.ModePerm to 0750(User=All, Group=Read/Execute, Others=None).
								if err := root.MkdirAll(filePath, 0750); err != nil {
									return fmt.Errorf("failed to crreate empty dir: %v", err)
								}
								continue
							}

							//file within dir
							if err := root.MkdirAll(filepath.Dir(filePath), 0750); err != nil {
								return fmt.Errorf("failed to unzip :%v", err)
							}

							//read-write, create, trucate config
							destFile, err := root.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
							if err != nil {
								return fmt.Errorf("failed to create empty dest: %v", err)
							}
							defer destFile.Close()

							//open file and copy contents
							fileInArchive, err := f.Open()
							if err != nil {
								return fmt.Errorf("failed to open file:%v", err)
							}
							defer fileInArchive.Close()

							const MxDecompress = 500 * 1024 * 1024 //500MiB

							if _, err := io.CopyN(destFile, fileInArchive, MxDecompress); err != nil && err != io.EOF {
								return fmt.Errorf("failed to copy contents or file too large: %v", err)
							}
							

						}
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
			// log.Fatal(err)
			fmt.Printf("GoSh error:%v\n", err)
		}

	}

}
