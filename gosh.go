package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	//  "os/exec"
	// "path/filepath"
	//"net/http"
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
	To get started: cli man
	  `)
	for scanner.Scan() {

		root := &cli.Command{

			//basic commands
			Commands: []*cli.Command{
				{
					Name:  "boom",
					Usage: "make an explosive entrance",
					Action: func(context.Context, *cli.Command) error {
						fmt.Println("boom! I say!")
						return nil
					},
				},
				{
					Name:  "man",
					Usage: "user manual",
					Action: func(context.Context, *cli.Command) error {
						fmt.Println(`	cat  | read contents
 	dir  | dir contents
	echo | Display Text
 	cd   | Change Directory
	pwd  | Print Working Directory
 	exit | exit the shell
	ls   | List Directory Contents
	mkdir| Make Directories
	rm   | Remove Files or Directories
	touch| change file timestamps or create an empty file
	mv   | move/rename file
	dir  | directory contents
	cat  | read contents
	info | file info
	
	for the rest of the commands, man {command} --help
									 `)
						return nil
					},
				},
				{
					Name:  "echo",
					Usage: "Display Text",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "help",
							Aliases: []string{"h","help"},
							Usage: "Display Text",
						},
						
					},
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println(c.Args().Get(0))
						if c.Bool("help"){
							fmt.Println("echo hello")
						}
						return nil
					},
				},
				{
					Name:  "cd",
					Usage: "Change Directory",
					//cli cd example
					Action: func(ctx context.Context, c *cli.Command) error {
						os.Chdir(c.Args().Get(0))
						fmt.Println(os.Getwd())
						return nil
					},
				},
				{
					Name:  "pwd",
					Usage: "Print Working Directory",
					// cli pwd
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println(os.Getwd())
						return nil
					},
				},
				{
					Name:  "exit",
					Usage: "exit",
					////cli exit
					Action: func(context.Context, *cli.Command) error {
						os.Exit(0)
						return nil
					},
				},
				{
					Name:  "ls",
					Usage: "List Directory Contents",
					//cli ls
					Action: func(context.Context, *cli.Command) error {
						fmt.Println(os.ReadDir("./"))
						return nil
					},
				},
				{
					Name:  "mkdir",
					Usage: "Make Directories",
					//cli mkdir example
					Action: func(ctx context.Context, c *cli.Command) error {
						os.Mkdir(c.Args().Get(0), 64)
						return nil
					},
				},
				{
					Name:  "rm",
					Usage: "Remove Files or Directories",
					//cli rm example.go
					//cli rm rf example
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "rf",
							Aliases: []string{"rf","r"},
							Usage: "recursive delete",
						},
						
					},
					Action: func(ctx context.Context, c *cli.Command) error {
						if c.Args().Len()==0 {
							// return cli.Exit("no file specified",14)
							fmt.Println("No file specified")
							return nil
						}
						if c.Bool("rf") {
							fmt.Println("are you sure you want to trigger recursive deletion ? y/n")
						
							response:=""
							fmt.Scanln(&response)
							if (response == "y" || response== "Y"){
								os.RemoveAll(c.Args().Get(0))
							}else{fmt.Println("Aborted")}
							
						}else{
							// for i:=0 ;i<len(c.Args().Slice());i++ {
							// 		os.Remove(c.Args().Get(i))
							// }
							for _,filename:=range c.Args().Slice() {
								err := os.Remove(filename)
								if err!=nil{
									fmt.Printf("Error deleting %s: %v\n",filename,err)
								}
							}
						}
						return nil
					},
					
				},

				{
					Name:  "touch",
					Usage: "change file timestamps or create an empty file",
					//cli touch example.go
					Action: func(ctx context.Context, c *cli.Command) error {
						os.Create(c.Args().Get(0))
						return nil
					},
				},
				{
					Name:  "mv",
					Usage: "move/rename file",
					//cli rn example exmp
					Action: func(ctx context.Context, c *cli.Command) error {
						os.Rename(c.Args().Get(0), c.Args().Get(1))
						return nil
					},
				},
				{
					Name:  "dir",
					Usage: "directory contents",
					//cli dir ./
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println(os.ReadDir(c.Args().Get(0)))
						return nil
					},
				},
				{
					Name:  "cat",
					Usage: "read contents",
					Action: func(ctx context.Context, c *cli.Command) error {
						data,err:=os.ReadFile(c.Args().Get(0))
						if err!=nil {
							log.Fatal(err)
						}
						os.Stdout.Write(data)
						return nil
					},
				},
				{
					Name:  "info",
					Usage: "file info",
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println(os.Stat(c.Args().Get(0)))
						return nil
					},
				},

				

				//////////////////////////////////////////////////////////////////////////////
				//system monitoring
				{
					Name:  "ps",
					Usage: "process status",
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println(os.Getpid())
						return nil
					},
				},
				{
					Name:  "uptime",
					Usage: "System Runtime",
				},
				{
					Name:  "sys",
					Usage: "System info",
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println("number of available cpu:", runtime.NumCPU(), "\n", "go version:", runtime.Version())
						return nil
					},
				},
				{
					Name:  "free",
					Usage: "Display Free and Used Memory",
					Action: func(ctx context.Context, c *cli.Command) error {
						fmt.Println()
						return nil
					},
				},
				// {
				// 	Name:  "kill",
				// 	Usage: "Terminate Processes",
				// 	Action: func( ctx context.Context, c * cli.Command)error{
				// 		os.Kill(c.Args().Get(0))
				// 		return nil
				// 	},

				// },

			},

				//////////////////////////////////////////////////////////////////////////////
				//Text Processing

				//////////////////////////////////////////////////////////////////////////////
				//Networking



				
		Flags: []cli.Flag{
            &cli.BoolFlag{
                Name:  "ginger-crouton",
                Usage: "is it in the soup?",
            },
        },
        // Action: func(ctx context.Context, cmd *cli.Command) error {
        //     if !cmd.Bool("ginger-crouton") {
        //         return cli.Exit("Ginger croutons are not in the soup", 86)
        //     }
        //     return nil
        // },
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
