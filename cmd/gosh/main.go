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
	args := os.Args[1:]

	root := &cli.Command{
		Name:    "gosh",
		Usage:   "GoSh - interactive shell",
		Commands: handlers.All(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}

	if len(args) > 0 {
		argsWithPrefix := append([]string{"gosh"}, args...)
		err := root.Run(context.Background(), argsWithPrefix)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	printBanner()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			fmt.Print("gosh> ")
			continue
		}
		fields := strings.Fields(text)
		// Strip "gosh" prefix if present
		if len(fields) > 0 && fields[0] == "gosh" {
			fields = fields[1:]
		}
		if len(fields) == 0 {
			fmt.Print("gosh> ")
			continue
		}
		// Add prefix for urfave routing
		argsWithPrefix := append([]string{"gosh"}, fields...)
		if err := root.Run(context.Background(), argsWithPrefix); err != nil {
			fmt.Printf("error: %v\n", err)
		}
		fmt.Print("gosh> ")
	}
}

func printBanner() {
	fmt.Print("\033[H\033[2J")
	fmt.Println(`
       ^  ^  ^   ^      ___I_      ^  ^   ^  ^  ^   ^  ^
      /|\/|\/|\ /|\    /\-_--\    /|\/|\ /|\/|\/|\ /|\/|\
      /|\/|\/|\ /|\   /  \_-__\   /|\/|\ /|\/|\/|\ /|\/|\
      /|\/|\/|\ /|\   |[]| [] |   /|\/|\ /|\/|\/|\ /|\/|\

	Welcome to GoSh v2.0.0
	To get started: {command} -h or {command}
	  `)
	fmt.Print("gosh> ")
}