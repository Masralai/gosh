package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Masralai/gosh/internal/handlers"
	"github.com/urfave/cli/v3"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\033[H\033[2J")
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
			Commands:  handlers.All(),
			Flags:     []cli.Flag{},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return nil
			},
		}
		if err := root.Run(context.Background(), append([]string{"GoSh"}, strings.Fields(scanner.Text())...)); err != nil {
			fmt.Printf("GoSh error: %v\n", err)
		}
	}
}