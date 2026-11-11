package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/urfave/cli/v3"
)

func init() {
	// shortcut: free -h for ls human-readable flag; help remains available via builtin `help`
	cli.HelpFlag = &cli.BoolFlag{Name: "help", Usage: "show help"}
}

func goshHistoryFilePath() string {
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

func Help() *cli.Command {
	return &cli.Command{
		Name:      "help",
		Usage:     "Show help for built-in commands",
		UsageText: "help [command]",
		Action: func(ctx context.Context, c *cli.Command) error {
			cmds := All()
			sort.Slice(cmds, func(i, j int) bool { return cmds[i].Name < cmds[j].Name })
			if c.Args().Len() == 0 {
				for _, cmd := range cmds {
					fmt.Printf("%-10s %s\n", cmd.Name, cmd.Usage)
				}
				return nil
			}
			target := c.Args().Get(0)
			for _, cmd := range cmds {
				if cmd.Name == target {
					printCommandHelp(cmd)
					return nil
				}
				for _, a := range cmd.Aliases {
					if a == target {
						printCommandHelp(cmd)
						return nil
					}
				}
			}
			return fmt.Errorf("unknown help topic: %s", target)
		},
	}
}

func printCommandHelp(cmd *cli.Command) {
	fmt.Printf("Name: %s\n", cmd.Name)
	if len(cmd.Aliases) > 0 {
		fmt.Printf("Aliases: %s\n", strings.Join(cmd.Aliases, ", "))
	}
	if cmd.Usage != "" {
		fmt.Printf("Usage: %s\n", cmd.Usage)
	}
	if cmd.UsageText != "" {
		fmt.Printf("UsageText: %s\n", cmd.UsageText)
	}
	if len(cmd.Flags) > 0 {
		fmt.Println("Flags:")
		for _, f := range cmd.Flags {
			names := f.Names()
			if len(names) == 0 {
				continue
			}
			fmt.Printf("  -%s  %s\n", strings.Join(names, ", -"), fmt.Sprint(f))
		}
	}
}

func History() *cli.Command {
	return &cli.Command{
		Name:      "history",
		Usage:     "Display or clear command history",
		UsageText: "history [-c]",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "c", Usage: "Clear the history"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			path := goshHistoryFilePath()
			if path == "" {
				return fmt.Errorf("history file path not available")
			}
			if c.Bool("c") {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					return fmt.Errorf("failed to prepare history file: %v", err)
				}
				if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
					return fmt.Errorf("failed to clear history: %v", err)
				}
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return fmt.Errorf("failed to read history: %v", err)
			}
			scanner := bufio.NewScanner(strings.NewReader(string(data)))
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			return nil
		},
	}
}

func Env() *cli.Command {
	return &cli.Command{
		Name:      "env",
		Usage:     "Display environment variables",
		UsageText: "env [name...]",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				envs := os.Environ()
				sort.Strings(envs)
				for _, e := range envs {
					fmt.Println(e)
				}
				return nil
			}
			for _, name := range c.Args().Slice() {
				fmt.Println(os.Getenv(name))
			}
			return nil
		},
	}
}

func Which() *cli.Command {
	return &cli.Command{
		Name:      "which",
		Usage:     "Locate a command",
		UsageText: "which <name> [name...]",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: which <name> [name...]")
			}
			cmds := All()
			isBuiltin := func(name string) (string, bool) {
				for _, cmd := range cmds {
					if cmd.Name == name {
						return cmd.Name, true
					}
					for _, a := range cmd.Aliases {
						if a == name {
							return cmd.Name + " (alias)", true
						}
					}
				}
				return "", false
			}
			var retErr error
			for _, name := range c.Args().Slice() {
				if resolved, ok := isBuiltin(name); ok {
					fmt.Printf("%s: gosh builtin %s\n", name, resolved)
					continue
				}
				if p, err := exec.LookPath(name); err == nil {
					fmt.Println(p)
					continue
				}
				fmt.Printf("%s: not found\n", name)
				retErr = fmt.Errorf("%s: not found", name)
			}
			return retErr
		},
	}
}

func Find() *cli.Command {
	return &cli.Command{
		Name:      "find",
		Usage:     "Search for files by name pattern",
		UsageText: "find <path> -name <pattern>",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Usage: "Base name pattern to match"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			path := "."
			args := c.Args().Slice()
			if len(args) > 0 && args[0] != "-name" {
				path = args[0]
			}
			pattern := c.String("name")
			if pattern == "" {
				// try to parse manual -name for flexibility
				for i := 0; i < len(args)-1; i++ {
					if args[i] == "-name" {
						pattern = args[i+1]
					}
				}
			}
			if pattern == "" {
				return fmt.Errorf("usage: find <path> -name <pattern>")
			}
			info, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("find: %s: %v", path, err)
			}
			if !info.IsDir() {
				if ok, _ := filepath.Match(pattern, filepath.Base(path)); ok {
					fmt.Println(filepath.ToSlash(path))
				}
				return nil
			}
			return filepath.WalkDir(path, func(p string, d os.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if ok, _ := filepath.Match(pattern, d.Name()); ok {
					fmt.Println(filepath.ToSlash(p))
				}
				return nil
			})
		},
	}
}
