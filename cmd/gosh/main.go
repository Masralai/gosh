package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
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

	histFile := historyFilePath()
	if histFile != "" {
		_ = os.MkdirAll(filepath.Dir(histFile), 0o750)
	}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       prompt(),
		HistoryFile:  histFile,
		AutoComplete: builtinCompleter(root),
	})
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
		args, perr := parseCommandLine(text)
		if perr != nil {
			fmt.Println("error:", perr)
			continue
		}
		runArgs(root, args, text)
	}
}

func run(root *cli.Command, args []string) {
	runArgs(root, args, strings.Join(args, " "))
}

func runArgs(root *cli.Command, args []string, line string) {
	if len(args) == 0 {
		return
	}
	if args[0] == "command" {
		if len(args) == 1 {
			fmt.Println("usage: command <cmd>")
			return
		}
		if strings.TrimSpace(line) != "" {
			forward(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "command")))
		} else {
			forward(strings.Join(args[1:], " "))
		}
		return
	}
	if isBuiltin(root, args[0]) {
		if err := root.Run(context.Background(), append([]string{"gosh"}, args...)); err != nil {
			fmt.Println("error:", err)
		}
		return
	}
	if strings.TrimSpace(line) != "" {
		forward(line)
	} else {
		forward(strings.Join(args, " "))
	}
}

func parseCommandLine(line string) ([]string, error) {
	args := []string{}
	var current strings.Builder
	inSingle, inDouble := false, false
	escaped := false
	tokenStarted := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if escaped {
			current.WriteByte(c)
			escaped = false
			tokenStarted = true
			continue
		}
		switch {
		case c == '\\' && !inSingle:
			if inDouble {
				if i+1 < len(line) {
					next := line[i+1]
					if next == '"' || next == '\\' {
						escaped = true
						tokenStarted = true
						continue
					}
				}
				current.WriteByte(c)
				tokenStarted = true
			} else {
				escaped = true
				tokenStarted = true
			}
		case c == '\'' && !inDouble:
			inSingle = !inSingle
			tokenStarted = true
		case c == '"' && !inSingle:
			inDouble = !inDouble
			tokenStarted = true
		case (c == ' ' || c == '\t' || c == '\n' || c == '\r') && !inSingle && !inDouble:
			if tokenStarted {
				args = append(args, current.String())
				current.Reset()
				tokenStarted = false
			}
		default:
			current.WriteByte(c)
			tokenStarted = true
		}
	}
	if escaped {
		current.WriteByte('\\')
	}
	if inSingle || inDouble {
		return nil, fmt.Errorf("unterminated quote")
	}
	if tokenStarted {
		args = append(args, current.String())
	}
	return args, nil
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

func historyFilePath() string {
	if p := os.Getenv("GOSH_HISTORY"); p != "" {
		return p
	}
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "gosh", "history")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".gosh_history")
	}
	return ""
}

func builtinCompletionNames(root *cli.Command) []string {
	seen := map[string]bool{}
	var names []string
	for _, c := range root.Commands {
		if c.Name != "" && !seen[c.Name] {
			seen[c.Name] = true
			names = append(names, c.Name)
		}
		for _, a := range c.Aliases {
			if a != "" && !seen[a] {
				seen[a] = true
				names = append(names, a)
			}
		}
	}
	sort.Strings(names)
	return names
}

func builtinCompleter(root *cli.Command) readline.AutoCompleter {
	items := builtinCompletionNames(root)
	pcs := make([]readline.PrefixCompleterInterface, 0, len(items))
	for _, name := range items {
		pcs = append(pcs, readline.PcItem(name))
	}
	return readline.NewPrefixCompleter(pcs...)
}
