package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

func Cli() *cli.Command {
	return &cli.Command{
		Name:            "cli",
		SkipFlagParsing: true,
		Usage:           "Echo the provided arguments",
		UsageText:       "cli [arguments]",
		Action: func(ctx context.Context, c *cli.Command) error {
			fmt.Println("cli", c.Args())
			return nil
		},
	}
}

func Boom() *cli.Command {
	return &cli.Command{
		Name:      "boom",
		Usage:     "Make an explosive entrance",
		UsageText: "boom",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}
}

func Echo() *cli.Command {
	return &cli.Command{
		Name:      "echo",
		Usage:     "Display a line of text",
		UsageText: "echo [-n] [-e] [text...]",
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
			text := strings.Join(c.Args().Slice(), " ")
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

func Clear() *cli.Command {
	return &cli.Command{
		Name:      "clear",
		Usage:     "Clear the terminal screen",
		UsageText: "clear",
		Action: func(ctx context.Context, c *cli.Command) error {
			if runtime.GOOS == "windows" {
				cmd := exec.Command("cmd", "/c", "cls")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				if err := cmd.Run(); err == nil {
					return nil
				}
			}
			fmt.Print("\033[H\033[2J")
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
			case 'a':
				result.WriteByte('\a')
				i++
				continue
			case 'b':
				result.WriteByte('\b')
				i++
				continue
			case 'f':
				result.WriteByte('\f')
				i++
				continue
			case 'r':
				result.WriteByte('\r')
				i++
				continue
			case 'v':
				result.WriteByte('\v')
				i++
				continue
			case 'x':
				if i+3 < len(s) {
					hex := s[i+2 : i+4]
					var n int
					_, err := fmt.Sscanf(hex, "%2x", &n)
					if err == nil && n >= 0 && n <= 255 {
						result.WriteByte(byte(n))
						i += 3
						continue
					}
				}
			}
		}
		result.WriteByte(s[i])
	}
	return result.String()
}

func Cd() *cli.Command {
	return &cli.Command{
		Name:      "cd",
		Usage:     "Change the current working directory",
		UsageText: "cd <path>",
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
		Name:      "pwd",
		Usage:     "Print the current working directory",
		UsageText: "pwd",
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
		Name:      "exit",
		Usage:     "Exit the shell",
		UsageText: "exit",
		Action: func(context.Context, *cli.Command) error {
			os.Exit(0)
			return nil
		},
	}
}

func Ls() *cli.Command {
	return &cli.Command{
		Name:                   "ls",
		Usage:                  "List directory contents",
		UsageText:              "ls [-a] [-l] [-h] [-R] [-S] [-t] [path]",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "a", Usage: "Include hidden files"},
			&cli.BoolFlag{Name: "l", Usage: "Use a long listing format"},
			&cli.BoolFlag{Name: "h", Usage: "With -l, print sizes in human readable format"},
			&cli.BoolFlag{Name: "R", Usage: "List subdirectories recursively"},
			&cli.BoolFlag{Name: "S", Usage: "Sort by file size"},
			&cli.BoolFlag{Name: "t", Usage: "Sort by modification time, newest first"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			path := "."
			if c.Args().Len() > 0 {
				path = c.Args().Get(0)
			}
			showAll := c.Bool("a")
			long := c.Bool("l")
			human := c.Bool("h")
			recursive := c.Bool("R")
			sortBySize := c.Bool("S")
			sortByTime := c.Bool("t")

			if !recursive {
				return listDir(path, showAll, long, human, sortBySize, sortByTime)
			}

			first := true
			return filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					return nil
				}
				if p != path && !showAll && strings.HasPrefix(d.Name(), ".") {
					return filepath.SkipDir
				}
				if !first {
					fmt.Println()
				}
				first = false
				fmt.Println(p + ":")
				if err := listDir(p, showAll, long, human, sortBySize, sortByTime); err != nil {
					return err
				}
				return nil
			})
		},
	}
}

type lsItem struct {
	entry os.DirEntry
	info  os.FileInfo
	name  string
}

func listDir(dir string, showAll, long, human, sortBySize, sortByTime bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to list directory contents: %v", err)
	}
	items := make([]lsItem, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if !showAll && strings.HasPrefix(name, ".") {
			continue
		}
		info, _ := e.Info()
		items = append(items, lsItem{entry: e, info: info, name: name})
	}
	sort.Slice(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if sortBySize {
			as, bs := int64(0), int64(0)
			if a.info != nil {
				as = a.info.Size()
			}
			if b.info != nil {
				bs = b.info.Size()
			}
			if as != bs {
				return as > bs
			}
		} else if sortByTime {
			at, bt := time.Time{}, time.Time{}
			if a.info != nil {
				at = a.info.ModTime()
			}
			if b.info != nil {
				bt = b.info.ModTime()
			}
			if !at.Equal(bt) {
				return at.After(bt)
			}
		}
		return a.name < b.name
	})
	for _, it := range items {
		if long {
			mode := "?????????"
			size := int64(0)
			mod := ""
			if it.info != nil {
				mode = it.info.Mode().String()
				size = it.info.Size()
				mod = it.info.ModTime().Format("Jan 02 15:04")
			}
			sizeStr := fmt.Sprintf("%d", size)
			if human {
				sizeStr = humanSize(size)
			}
			name := it.name
			if it.entry.IsDir() {
				name += "/"
			}
			fmt.Printf("%s %8s %s %s\n", mode, sizeStr, mod, name)
		} else {
			if it.entry.IsDir() {
				fmt.Printf("%s/\n", it.name)
			} else {
				fmt.Printf("%s\n", it.name)
			}
		}
	}
	return nil
}

func humanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	suffix := []string{"K", "M", "G", "T", "P"}
	if exp >= len(suffix) {
		exp = len(suffix) - 1
	}
	return fmt.Sprintf("%.1f%s", float64(n)/float64(div), suffix[exp])
}

func Mkdir() *cli.Command {
	return &cli.Command{
		Name:      "mkdir",
		Usage:     "Create directories",
		UsageText: "mkdir [-p] <path>",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "p", Usage: "Create parent directories as needed; no error if existing"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: mkdir [-p] <path>")
			}
			path := c.Args().Get(0)
			var err error
			if c.Bool("p") {
				err = os.MkdirAll(path, 0o750)
			} else {
				err = os.Mkdir(path, 0o750)
			}
			if err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			return nil
		},
	}
}

func Rm() *cli.Command {
	return &cli.Command{
		Name:                   "rm",
		Usage:                  "Remove files or directories",
		UsageText:              "rm [-r|-R] [-f] <path>...",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "r", Aliases: []string{"R"}, Usage: "Remove directories and their contents recursively"},
			&cli.BoolFlag{Name: "f", Usage: "Ignore nonexistent files and arguments, never prompt"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: rm [-r|-R] [-f] <path>...")
			}
			recursive := c.Bool("r") || c.Bool("R")
			force := c.Bool("f")
			for _, name := range c.Args().Slice() {
				info, err := os.Lstat(name)
				if err != nil {
					if os.IsNotExist(err) && force {
						continue
					}
					return fmt.Errorf("failed to remove %s: %v", name, err)
				}
				if info.IsDir() && !recursive {
					return fmt.Errorf("failed to remove %s: is a directory", name)
				}
				if recursive {
					if err := os.RemoveAll(name); err != nil {
						if force && os.IsNotExist(err) {
							continue
						}
						return fmt.Errorf("failed to remove %s: %v", name, err)
					}
				} else {
					if err := os.Remove(name); err != nil {
						if force && os.IsNotExist(err) {
							continue
						}
						return fmt.Errorf("failed to remove %s: %v", name, err)
					}
				}
			}
			return nil
		},
	}
}

func Touch() *cli.Command {
	return &cli.Command{
		Name:      "touch",
		Usage:     "Create an empty file or update its timestamp",
		UsageText: "touch <file> [file...]",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: touch <file> [file...]")
			}
			now := time.Now()
			for _, name := range c.Args().Slice() {
				_, err := os.Stat(name)
				if err != nil {
					if os.IsNotExist(err) {
						// #nosec G304,G302 -- touch creates user-specified file path
						f, cerr := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0o600)
						if cerr != nil {
							return fmt.Errorf("failed to create file: %v", cerr)
						}
						if cerr := f.Close(); cerr != nil {
							return fmt.Errorf("failed to close file: %v", cerr)
						}
					} else {
						return fmt.Errorf("failed to access file: %v", err)
					}
				}
				if err := os.Chtimes(name, now, now); err != nil {
					return fmt.Errorf("failed to update timestamp: %v", err)
				}
			}
			return nil
		},
	}
}

func Mv() *cli.Command {
	return &cli.Command{
		Name:      "mv",
		Usage:     "Move or rename a file or directory",
		UsageText: "mv <source> <destination>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("usage: mv <source> <destination>")
			}
			src := c.Args().Get(0)
			dest := c.Args().Get(1)
			if err := os.Rename(src, dest); err != nil {
				return fmt.Errorf("failed to move %q: %v", src, err)
			}
			return nil
		},
	}
}

func Cp() *cli.Command {
	return &cli.Command{
		Name:      "cp",
		Usage:     "Copy files and directories",
		UsageText: "cp [-r] <source> <destination>",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "r", Aliases: []string{"R"}, Usage: "Copy directories recursively"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("usage: cp [-r] <source> <destination>")
			}
			recursive := c.Bool("r") || c.Bool("R")
			src := filepath.Clean(c.Args().Get(0))
			dest := filepath.Clean(c.Args().Get(1))

			info, err := os.Stat(src)
			if err != nil {
				return fmt.Errorf("failed to access source: %v", err)
			}
			if info.IsDir() {
				if !recursive {
					return fmt.Errorf("cp: -r not specified; omitting directory '%s'", src)
				}
				return copyDir(src, dest)
			}
			return copyFile(src, dest, info.Mode().Perm())
		},
	}
}

func copyFile(src, dest string, perm os.FileMode) error {
	if dInfo, derr := os.Stat(dest); derr == nil && dInfo.IsDir() {
		dest = filepath.Join(dest, filepath.Base(src))
	}
	// #nosec G304 -- source/dest paths are user-specified cp operands
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %v", err)
	}
	defer srcFile.Close()

	// #nosec G304 -- destination path is user-specified cp operand
	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("failed to create destination: %v", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy: %v", err)
	}
	return nil
}

func copyDir(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to access source: %v", err)
	}
	if dInfo, derr := os.Stat(dest); derr == nil && dInfo.IsDir() {
		dest = filepath.Join(dest, filepath.Base(src))
	}
	if err := os.MkdirAll(dest, srcInfo.Mode().Perm()); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dest, err)
	}
	return filepath.WalkDir(src, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			info, ierr := d.Info()
			mode := srcInfo.Mode().Perm()
			if ierr == nil {
				mode = info.Mode().Perm()
			}
			return os.MkdirAll(target, mode)
		}
		info, ierr := d.Info()
		mode := os.FileMode(0o644)
		if ierr == nil {
			mode = info.Mode().Perm()
		}
		return copyFile(path, target, mode)
	})
}

func Dir() *cli.Command {
	return &cli.Command{
		Name:      "dir",
		Usage:     "Display directory contents",
		UsageText: "dir [path]",
		Action: func(ctx context.Context, c *cli.Command) error {
			path := "."
			if c.Args().Len() > 0 {
				path = c.Args().Get(0)
			}
			entries, err := os.ReadDir(path)
			if err != nil {
				return fmt.Errorf("failed to read directory contents: %v", err)
			}
			sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
			for _, e := range entries {
				name := e.Name()
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

func Cat() *cli.Command {
	return &cli.Command{
		Name:      "cat",
		Usage:     "Read and display file contents",
		UsageText: "cat [-n] [-b] [-s] <file> [file...]",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "n", Usage: "Number all output lines"},
			&cli.BoolFlag{Name: "b", Usage: "Number non-blank output lines (overrides -n)"},
			&cli.BoolFlag{Name: "s", Usage: "Squeeze multiple adjacent blank lines"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: cat [-n] [-b] [-s] <file> [file...]")
			}
			numberAll := c.Bool("n")
			numberNonBlank := c.Bool("b")
			squeeze := c.Bool("s")
			if numberNonBlank {
				numberAll = false
			}
			lineNum := 0
			prevBlank := false
			for _, name := range c.Args().Slice() {
				// #nosec G304 -- file paths are user-specified cat operands
				f, err := os.Open(name)
				if err != nil {
					return fmt.Errorf("failed to open file: %v", err)
				}
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					line := scanner.Text()
					isBlank := len(strings.TrimSpace(line)) == 0
					if squeeze && isBlank && prevBlank {
						continue
					}
					prevBlank = isBlank
					switch {
					case numberAll:
						lineNum++
						fmt.Printf("%6d  %s\n", lineNum, line)
					case numberNonBlank:
						if !isBlank {
							lineNum++
							fmt.Printf("%6d  %s\n", lineNum, line)
						} else {
							fmt.Println()
						}
					default:
						fmt.Println(line)
					}
				}
				err = scanner.Err()
				_ = f.Close()
				if err != nil {
					return fmt.Errorf("failed to read file: %v", err)
				}
			}
			return nil
		},
	}
}

func Info() *cli.Command {
	return &cli.Command{
		Name:      "info",
		Usage:     "Display file information",
		UsageText: "info <path>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: info <path>")
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
