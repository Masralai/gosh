package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Masralai/gosh/internal/handlers"
	"github.com/reeflective/readline"
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
			cmd := exec.Command("bash", "-c", bashCmd)
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

	rl := readline.NewShell()

	rl.Prompt.Primary(func() string { 
		return "gosh> "
	})

	rl.Config.Set("completion-auto", false)
	rl.Config.Set("show-all-if-ambiguous", false)
	rl.Config.Set("menu-complete-display-prefix", false)
	
	rl.Completer = func(line []rune, cursor int) readline.Completions {
		input := string(line)
		completions := handlers.GetCompletionsForLine(input)
		if len(completions) == 0 {
			return readline.Completions{}
		}
		return readline.CompleteValues(completions...)
	}

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		text := strings.TrimSpace(line)
		if text == "" {
			continue
		}

		fields := strings.Fields(text)
		if len(fields) > 0 && fields[0] == "gosh" {
			fields = fields[1:]
		}
		if len(fields) == 0 {
			continue
		}

		if strings.HasPrefix(fields[0], "!") {
			cmd := strings.TrimPrefix(fields[0], "!")
			argsWithPrefix := append([]string{"gosh", cmd}, fields[1:]...)
			if err := root.Run(context.Background(), argsWithPrefix); err != nil {
				fmt.Printf("error: %v\n", err)
			}
			handlers.GetHistoryManager().Add(text)
			continue
		}

		bashCmd := strings.Join(fields, " ")
		bashCmd = strings.TrimSpace(bashCmd)
		if bashCmd == "" {
			continue
		}
		cmd := exec.Command("bash", "-c", bashCmd)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("error: %v\n", err)
		}
		handlers.GetHistoryManager().Add(text)
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