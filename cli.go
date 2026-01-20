package main

import (
	"bufio"
	"strings"

	//"strings"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		root := &cli.Command{

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
					Name:  "cd",
					Usage: "ch dir",
					Action: func(ctx context.Context,c *cli.Command) error {
						os.Chdir(c.Args().Get(0))
                        fmt.Println(os.Getwd())
						return nil
					},
				},
                {
					Name:  "pwd",
					Usage: "print dir",
					Action: func(ctx context.Context,c *cli.Command) error {
						fmt.Println(os.Getwd())
						return nil
					},
				},
				{
					Name:  "exit",
					Usage: "exit",
					Action: func(context.Context, *cli.Command) error {
						os.Exit(0)
						return nil
					},
				},
                {
					Name:  "ls",
					Usage: "list",
					Action: func(context.Context, *cli.Command) error {
						fmt.Println(os.ReadDir("./"))
						return nil
					},
				},
                {
					Name:  "mkdir",
					Usage: "make dir",
					Action: func(ctx context.Context,c *cli.Command) error {
						os.Mkdir(c.Args().Get(0),64)
						return nil
					},
				},
                {
					Name:  "rmdir",
					Usage: "remove dir",
					Action: func(ctx context.Context,c *cli.Command) error {
						os.Remove(c.Args().Get(0))
						return nil
					},
				},
                {
					Name:  "rmall",
					Usage: "remove all",
					Action: func(ctx context.Context,c *cli.Command) error {
						os.RemoveAll(c.Args().Get(0))
						return nil
					},
				},
                {
					Name:  "touch",
					Usage: "init file",
					Action: func(ctx context.Context,c *cli.Command) error {
						os.Create(c.Args().Get(0))
						return nil
					},
				},
                {
					Name:  "rn",
					Usage: "rename dir",
					Action: func(ctx context.Context,c *cli.Command) error {
						os.Rename(c.Args().Get(0),c.Args().Get(1))
						return nil
					},
				},
                {
					Name:  "dir",
					Usage: "dir contents",
					Action: func(ctx context.Context,c *cli.Command) error {
						os.ReadDir(c.Args().Get(0))
						return nil
					},
				},
                {
					Name:  "cat",
					Usage: "read contents",
					Action: func(ctx context.Context,c *cli.Command) error {
						fmt.Println(os.ReadFile(scanner.Text()))
						return nil
					},
				},
			},
		}
		if err := root.Run(context.Background(), strings.Fields(scanner.Text())); err != nil {
			log.Fatal(err)
		}

	}

}
