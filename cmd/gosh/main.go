package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Masralai/gosh/internal/handlers"
	"github.com/chzyer/readline"
	"github.com/urfave/cli/v3"
)

func main() {
	args := os.Args[1:]

	root := &cli.Command{
		Name:     "gosh",
		Usage:    "GoSh - interactive shell",
		Commands: handlers.All(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}

	if len(args) > 0 {
		if len(args) > 0 && strings.HasPrefix(args[0], "!") {
			cmd := strings.TrimPrefix(args[0], "!")
			argsWithPrefix := append([]string{"gosh", cmd}, args[1:]...)
			err := root.Run(context.Background(), argsWithPrefix)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}
		} else {
			bashCmd := strings.Join(args, " ")
			cmd := exec.Command("bash", "-c", bashCmd) // #nosec G204 G702
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}
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
		_ = rl.SaveHistory(text)

		fields := strings.Fields(text)
		if len(fields) > 0 && fields[0] == "gosh" {
			fields = fields[1:]
		}
		if len(fields) == 0 {
			continue
		}

		// ! prefix for gosh commands, default to bash
		if strings.HasPrefix(fields[0], "!") {
			cmd := strings.TrimPrefix(fields[0], "!")
			argsWithPrefix := append([]string{"gosh", cmd}, fields[1:]...)
			if err := root.Run(context.Background(), argsWithPrefix); err != nil {
				fmt.Printf("error: %v\n", err)
			}
			continue
		}

		// Default: forward to bash
		bashCmd := strings.Join(fields, " ")
		bashCmd = strings.TrimSpace(bashCmd)
		if bashCmd == "" {
			continue
		}
// #nosec G204 G702
		cmd := exec.Command("bash", "-c", bashCmd)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
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

	Welcome to GoSh
	To get started: {command} -h or {command}
	  `)
}
