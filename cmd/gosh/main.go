package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Masralai/gosh/internal/handlers"
	"github.com/chzyer/readline"
	"github.com/urfave/cli/v3"
)

var (
	hostPath string
	hostFlag string
)

func main() {
	hostPath, hostFlag = findHostShell()
	root := &cli.Command{
		Name:     "gosh",
		Usage:    "GoSh - interactive shell",
		Commands: handlers.All(),
	}

	if args := os.Args[1:]; len(args) > 0 {
		run(root, args)
		return
	}

	printBanner()

	rl, err := readline.New(prompt())
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		rl.SetPrompt(prompt())
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
		run(root, strings.Fields(text))
	}
}

func run(root *cli.Command, args []string) {
	if len(args) == 0 {
		return
	}
	if args[0] == "command" {
		if len(args) == 1 {
			fmt.Println("usage: command <cmd>")
			return
		}
		forward(strings.Join(args[1:], " "))
		return
	}
	if isBuiltin(root, args[0]) {
		if err := root.Run(context.Background(), append([]string{"gosh"}, args...)); err != nil {
			fmt.Println("error:", err)
		}
		return
	}
	forward(strings.Join(args, " "))
}

func isBuiltin(root *cli.Command, name string) bool {
	for _, c := range root.Commands {
		if c.Name == name {
			return true
		}
		for _, a := range c.Aliases {
			if a == name {
				return true
			}
		}
	}
	return false
}

// findHostShell picks the shell that runs anything GoSh does not implement.
// shortcut: GOSH_SHELL is assumed to accept -c; cmd.exe is the only /C default.
func findHostShell() (string, string) {
	if s := os.Getenv("GOSH_SHELL"); s != "" {
		return s, "-c"
	}
	if p, err := exec.LookPath("bash"); err == nil && exec.Command(p, "-c", "exit 0").Run() == nil { // #nosec G204 - fixed argument, probed binary from LookPath
		return p, "-c"
	}
	if p, err := exec.LookPath("sh"); err == nil {
		return p, "-c"
	}
	if runtime.GOOS == "windows" {
		if p, err := exec.LookPath("cmd"); err == nil {
			return p, "/C"
		}
	}
	return "", ""
}

func forward(line string) {
	if hostPath == "" {
		fmt.Println("error: no host shell found (set GOSH_SHELL)")
		return
	}
	cmd := exec.Command(hostPath, hostFlag, line) // #nosec G204 G702 - running user input is the point
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// A non-zero exit is the command's result, not a GoSh failure.
		if _, ok := err.(*exec.ExitError); !ok {
			fmt.Println("error:", err)
		}
	}
}

func prompt() string {
	wd, err := os.Getwd()
	if err != nil {
		return "gosh ❯ "
	}
	home, _ := os.UserHomeDir()
	return shortPath(wd, home) + " ❯ "
}

func shortPath(wd, home string) string {
	p := filepath.ToSlash(wd)
	if home != "" {
		rel, err := filepath.Rel(home, wd)
		if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			if rel == "." {
				p = "~"
			} else {
				p = "~/" + filepath.ToSlash(rel)
			}
		}
	}
	if len(p) > 24 {
		parts := strings.Split(strings.TrimPrefix(p, "/"), "/")
		if len(parts) >= 2 {
			p = ".../" + strings.Join(parts[len(parts)-2:], "/")
		}
	}
	return p
}

func printBanner() {
	fmt.Print("\033[H\033[2J")
	fmt.Println(`
       ^  ^  ^   ^      ___I_      ^  ^   ^  ^  ^   ^  ^
      /|\/|\/|\ /|\    /\-_--\    /|\/|\ /|\/|\/|\ /|\/|\
      /|\/|\/|\ /|\   /  \_-__\   /|\/|\ /|\/|\/|\ /|\/|\
      /|\/|\/|\ /|\   |[]| [] |   /|\/|\ /|\/|\/|\ /|\/|\

	Welcome to GoSh
	Type a command, or command <cmd> for the host shell.
	  `)
}
