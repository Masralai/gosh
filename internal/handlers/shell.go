package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

func Cli() *cli.Command {
	return &cli.Command{
		Name:             "cli",
		SkipFlagParsing:   true,
		Usage:            "Echo the provided arguments",
		UsageText:        "cli [arguments]",
		Action: func(ctx context.Context, c *cli.Command) error {
			fmt.Println("cli", c.Args())
			return nil
		},
	}
}

func Boom() *cli.Command {
	return &cli.Command{
		Name:     "boom",
		Usage:    "Make an explosive entrance",
		UsageText: "cli boom",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}
}

func Echo() *cli.Command {
	return &cli.Command{
		Name:     "echo",
		Usage:    "Display text",
		UsageText: "cli echo [-n] [-e] <text>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "n",
				Usage: "Do not output the trailing newline",
			},
			&cli.BoolFlag{
				Name:  "e",
				Usage: "Enable interpretation of backslash escapes",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: echo <text>")
			}
			text := c.Args().Get(0)
			if c.Bool("e") {
				text = expandEscapes(text)
			}
			if c.Bool("n") {
				fmt.Print(text)
			} else {
				fmt.Println(text)
			}
			return nil
		},
	}
}

func expandEscapes(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result.WriteByte('\n')
				i++
				continue
			case 't':
				result.WriteByte('\t')
				i++
				continue
			case '\\':
				result.WriteByte('\\')
				i++
				continue
			}
		}
		result.WriteByte(s[i])
	}
	return result.String()
}

func Cd() *cli.Command {
	return &cli.Command{
		Name:     "cd",
		Usage:    "Change the current working directory",
		UsageText: "cli cd <path>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: cd <path>")
			}
			if err := os.Chdir(c.Args().Get(0)); err != nil {
				return fmt.Errorf("failed to change directory: %v", err)
			}
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %v", err)
			}
			fmt.Println(wd)
			return nil
		},
	}
}

func Pwd() *cli.Command {
	return &cli.Command{
		Name:     "pwd",
		Usage:    "Print the current working directory",
		UsageText: "cli pwd",
		Action: func(ctx context.Context, c *cli.Command) error {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %v", err)
			}
			fmt.Println(wd)
			return nil
		},
	}
}

func Exit() *cli.Command {
	return &cli.Command{
		Name:     "exit",
		Usage:    "Exit the shell",
		UsageText: "cli exit",
		Action: func(context.Context, *cli.Command) error {
			os.Exit(0)
			return nil
		},
	}
}

func Ls() *cli.Command {
	return &cli.Command{
		Name:     "ls",
		Usage:    "List directory contents",
		UsageText: "cli ls [-R] [-S] [-a] [path]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "R",
				Usage: "List subdirectories recursively",
			},
			&cli.BoolFlag{
				Name:  "S",
				Usage: "Sort by file size",
			},
			&cli.BoolFlag{
				Name:  "a",
				Usage: "Include hidden files",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			path := "."
			if c.Args().Len() > 0 {
				path = c.Args().Get(0)
			}
			entries, err := os.ReadDir(path)
			if err != nil {
				return fmt.Errorf("failed to list directory contents: %v", err)
			}
			for _, e := range entries {
				name := e.Name()
				if !c.Bool("a") && strings.HasPrefix(name, ".") {
					continue
				}
				if e.IsDir() {
					fmt.Printf("%s/\n", name)
				} else {
					fmt.Printf("%s\n", name)
				}
			}
			return nil
		},
	}
}

func Mkdir() *cli.Command {
	return &cli.Command{
		Name:     "mkdir",
		Usage:    "Create directories",
		UsageText: "cli mkdir <path>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: mkdir <pathname>")
			}
			if err := os.Mkdir(c.Args().Get(0), 0o755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			fmt.Printf("directory created: %s\n", c.Args().Get(0))
			return nil
		},
	}
}

func Rm() *cli.Command {
	return &cli.Command{
		Name:     "rm",
		Usage:    "Remove files or directories",
		UsageText: "cli rm  | cli rm -rf <path>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "rf",
				Aliases: []string{"r"},
				Usage:   "Remove directories and their contents recursively",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: rm  | rm -rf <path>")
			}
			if c.Bool("rf") {
				fmt.Println("are you sure you want to trigger recursive deletion? y/n")
				var response string
				if _, err := fmt.Scanln(&response); err != nil {
					return fmt.Errorf("failed to read input: %v", err)
				}
				if response == "y" || response == "Y" {
					if err := os.RemoveAll(c.Args().Get(0)); err != nil {
						return fmt.Errorf("failed recursive delete: %v", err)
					}
				} else {
					fmt.Println("Aborted")
				}
			} else {
				for _, name := range c.Args().Slice() {
					if err := os.Remove(name); err != nil {
						return fmt.Errorf("failed to delete %s: %v", name, err)
					}
				}
			}
			return nil
		},
	}
}

func Touch() *cli.Command {
	return &cli.Command{
		Name:     "touch",
		Usage:    "Create an empty file or update its timestamp",
		UsageText: "cli touch ",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: touch ")
			}
			f, err := os.Create(c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			f.Close()
			fmt.Printf("file created: %s\n", c.Args().Get(0))
			return nil
		},
	}
}

func Mv() *cli.Command {
	return &cli.Command{
		Name:     "mv",
		Usage:    "Move or rename a file or directory",
		UsageText: "cli mv <source> <destination>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("usage: mv <source> <destination>")
			}
			src := c.Args().Get(0)
			dest := c.Args().Get(1)
			if err := os.Rename(src, dest); err != nil {
				return fmt.Errorf("failed to move %q: %v", src, err)
			}
			fmt.Printf("moved %q to %q\n", src, dest)
			return nil
		},
	}
}

func Cp() *cli.Command {
	return &cli.Command{
		Name:     "cp",
		Usage:    "Copy files and directories",
		UsageText: "cli cp <source> <destination>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("usage: cp <source> <destination>")
			}
			src := c.Args().Get(0)
			dest := c.Args().Get(1)
			srcFS := os.DirFS(src)
			if err := os.CopyFS(dest, srcFS); err != nil {
				return fmt.Errorf("failed to copy %q: %v", src, err)
			}
			fmt.Printf("copied %q to %q\n", src, dest)
			return nil
		},
	}
}

func Dir() *cli.Command {
	return &cli.Command{
		Name:     "dir",
		Usage:    "Display directory contents",
		UsageText: "cli dir <path>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: dir <pathname>")
			}
			entries, err := os.ReadDir(c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to read directory contents: %v", err)
			}
			fmt.Println(entries)
			return nil
		},
	}
}

func Cat() *cli.Command {
	return &cli.Command{
		Name:     "cat",
		Usage:    "Read and display file contents",
		UsageText: "cli cat [-n] [-b] [-s] ",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "n",
				Usage: "Number all output lines",
			},
			&cli.BoolFlag{
				Name:  "b",
				Usage: "Number non-blank output lines",
			},
			&cli.BoolFlag{
				Name:  "s",
				Usage: "Squeeze multiple blank lines into one",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: cat ")
			}
			f, err := os.Open(c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			lineNum := 0
			prevEmpty := false
			for scanner.Scan() {
				line := scanner.Text()
				isEmpty := len(strings.TrimSpace(line)) == 0
				if c.Bool("s") && isEmpty && prevEmpty {
					continue
				}
				prevEmpty = isEmpty

				lineNum++
				if c.Bool("n") {
					fmt.Printf("%6d  %s\n", lineNum, line)
				} else if c.Bool("b") && !isEmpty {
					fmt.Printf("%6d  %s\n", lineNum, line)
				} else if !c.Bool("n") && !c.Bool("b") {
					fmt.Println(line)
				}
			}
			return nil
		},
	}
}

func Info() *cli.Command {
	return &cli.Command{
		Name:     "info",
		Usage:    "Display file information",
		UsageText: "cli info ",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: info ")
			}
			s, err := os.Stat(c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to get file info: %v", err)
			}
			fmt.Printf("File: %s\n", s.Name())
			fmt.Printf("Size: %d bytes\n", s.Size())
			fmt.Printf("Mode: %s\n", s.Mode())
			fmt.Printf("Last Modified: %s\n", s.ModTime().Format("2006-01-02 15:04:05"))
			fmt.Printf("Directory?: %v\n", s.IsDir())
			return nil
		},
	}
}