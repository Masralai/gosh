package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
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

	rl, err := readline.New("gosh> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		text := strings.TrimSpace(line)
		if text == "" {
			continue
		}
		// #nosec G104 - ignore history save errors
		rl.SaveHistory(text)

		fields := strings.Fields(text)
		if len(fields) > 0 && fields[0] == "gosh" {
			fields = fields[1:]
		}
		if len(fields) == 0 {
			continue
		}

		// Bash forwarding with ! prefix
		if strings.HasPrefix(fields[0], "!") {
			bashCmd := strings.TrimPrefix(fields[0], "!")
			if len(fields) > 1 {
				bashCmd += " " + strings.Join(fields[1:], " ")
			}
			bashCmd = strings.TrimSpace(bashCmd)
			if bashCmd == "" {
				fmt.Printf("error: empty bash command\n")
				continue
			}
			// #nosec G204 - explicit user command with ! prefix
			cmd := exec.Command("bash", "-c", bashCmd)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("error: %v\n", err)
			}
			continue
		}

		// Add prefix for urfave routing
		argsWithPrefix := append([]string{"gosh"}, fields...)
		if err := root.Run(context.Background(), argsWithPrefix); err != nil {
			fmt.Printf("error: %v\n", err)
		}
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
}