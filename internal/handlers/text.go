package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/urfave/cli/v3"
)

func Grep() *cli.Command {
	return &cli.Command{
		Name:      "grep",
		Usage:     "Search for a pattern in files",
		UsageText: "grep [-i] [-n] [-r] [-v] <pattern> <path>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "i",
				Aliases: []string{"f"}, // shortcut: -f kept as ignore-case compat alias, not pattern-file
				Usage:   "Ignore case differences",
			},
			&cli.BoolFlag{
				Name:  "n",
				Usage: "Show line numbers",
			},
			&cli.BoolFlag{
				Name:  "r",
				Usage: "Search directories recursively",
			},
			&cli.BoolFlag{
				Name:  "v",
				Usage: "Select non-matching lines",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("usage: grep [-i] [-n] [-r] [-v] <pattern> <path>")
			}

			filename := c.Args().Get(1)
			pattern := c.Args().Get(0)
			ignoreCase := c.Bool("i") || c.Bool("f")
			showNum := c.Bool("n")
			invert := c.Bool("v")

			var regObj *regexp.Regexp
			var err error
			if ignoreCase {
				regObj, err = regexp.Compile("(?i)" + pattern)
			} else {
				regObj, err = regexp.Compile(pattern)
			}
			if err != nil {
				return fmt.Errorf("failed to compile pattern: %v", err)
			}

			printMatch := func(prefix string, lineNum int, line string) {
				if showNum {
					if prefix != "" {
						fmt.Printf("%s:%d:%s\n", prefix, lineNum, line)
					} else {
						fmt.Printf("%d:%s\n", lineNum, line)
					}
				} else {
					if prefix != "" {
						fmt.Printf("%s:%s\n", prefix, line)
					} else {
						fmt.Println(line)
					}
				}
			}

			if c.Bool("r") {
				return filepath.WalkDir(filename, func(path string, d fs.DirEntry, walkErr error) error {
					if walkErr != nil {
						return walkErr
					}
					if d.IsDir() {
						return nil
					}
					// #nosec G304 G122
					file, err := os.Open(path)
					if err != nil {
						return nil
					}
					scanner := bufio.NewScanner(file)
					scanner.Buffer(make([]byte, 64*1024), 1024*1024)
					lineNum := 0
					for scanner.Scan() {
						lineNum++
						line := scanner.Text()
						matched := regObj.MatchString(line)
						if invert {
							matched = !matched
						}
						if matched {
							printMatch(path, lineNum, line)
						}
					}
					scanErr := scanner.Err()
					_ = file.Close()
					if scanErr != nil {
						return nil
					}
					return nil
				})
			}

			// #nosec G304 -- file path is user-specified grep target
			file, err := os.Open(filename)
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Buffer(make([]byte, 64*1024), 1024*1024)
			lineNum := 0
			for scanner.Scan() {
				lineNum++
				line := scanner.Text()
				matched := regObj.MatchString(line)
				if invert {
					matched = !matched
				}
				if matched {
					printMatch("", lineNum, line)
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading file: %v", err)
			}

			return nil
		},
	}
}

func Head() *cli.Command {
	return &cli.Command{
		Name:      "head",
		Usage:     "Display the beginning of a file",
		UsageText: "head [-n <lines>] <file>",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "n",
				Usage: "Number of lines to display",
				Value: 10,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: head [-n <lines>] <file>")
			}
			// #nosec G304 -- file path is user-specified head target
			file, err := os.Open(c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			count := 0
			maxLines := int(c.Int("n"))
			for scanner.Scan() {
				if count >= maxLines {
					break
				}
				fmt.Println(scanner.Text())
				count++
			}
			return nil
		},
	}
}

func Tail() *cli.Command {
	return &cli.Command{
		Name:      "tail",
		Usage:     "Display the end of a file",
		UsageText: "tail [-n <lines>] <file>",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "lines",
				Aliases: []string{"n"},
				Usage:   "Number of lines to display",
				Value:   10,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: tail [-n <lines>] <file>")
			}
			// #nosec G304 -- file path is user-specified tail target
			file, err := os.Open(c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}
			defer file.Close()

			var lines []string
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			n := int(c.Int("lines"))
			if n > len(lines) {
				n = len(lines)
			}
			start := len(lines) - n
			for i := start; i < len(lines); i++ {
				fmt.Println(lines[i])
			}
			return nil
		},
	}
}

func Wc() *cli.Command {
	return &cli.Command{
		Name:      "wc",
		Usage:     "Count lines, words, and bytes in files",
		UsageText: "wc [-l] [-w] [-c] <file> [file...]",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "l", Usage: "Count lines"},
			&cli.BoolFlag{Name: "w", Usage: "Count words"},
			&cli.BoolFlag{Name: "c", Usage: "Count bytes"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: wc [-l] [-w] [-c] <file> [file...]")
			}
			showL, showW, showC := c.Bool("l"), c.Bool("w"), c.Bool("c")
			if !showL && !showW && !showC {
				showL, showW, showC = true, true, true
			}
			var totL, totW, totC int
			for _, name := range c.Args().Slice() {
				l, w, b, err := wcFile(name)
				if err != nil {
					return fmt.Errorf("wc: %s: %v", name, err)
				}
				totL += l
				totW += w
				totC += b
				printWc(showL, showW, showC, l, w, b, name)
			}
			if c.Args().Len() > 1 {
				printWc(showL, showW, showC, totL, totW, totC, "total")
			}
			return nil
		},
	}
}

func wcFile(name string) (int, int, int, error) {
	// #nosec G304 -- file path is user-specified wc target
	f, err := os.Open(name)
	if err != nil {
		return 0, 0, 0, err
	}
	defer f.Close()
	var lines, words, bytesCount int
	inWord := false
	r := bufio.NewReader(f)
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, 0, err
		}
		bytesCount++
		if b == '\n' {
			lines++
		}
		if b == ' ' || b == '\n' || b == '\t' || b == '\r' || b == '\v' || b == '\f' {
			inWord = false
		} else if !inWord {
			inWord = true
			words++
		}
	}
	return lines, words, bytesCount, nil
}

func printWc(showL, showW, showC bool, l, w, b int, name string) {
	parts := []string{}
	if showL {
		parts = append(parts, fmt.Sprintf("%d", l))
	}
	if showW {
		parts = append(parts, fmt.Sprintf("%d", w))
	}
	if showC {
		parts = append(parts, fmt.Sprintf("%d", b))
	}
	fmt.Printf("%s %s\n", strings.Join(parts, " "), name)
}

func SortCmd() *cli.Command {
	return &cli.Command{
		Name:      "sort",
		Usage:     "Sort lines of text",
		UsageText: "sort [-r] <file>",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "r", Usage: "Reverse the result of comparisons"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: sort [-r] <file>")
			}
			// #nosec G304 -- file path is user-specified sort target
			file, err := os.Open(c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}
			defer file.Close()
			var lines []string
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to read file: %v", err)
			}
			if c.Bool("r") {
				sort.Slice(lines, func(i, j int) bool { return lines[i] > lines[j] })
			} else {
				sort.Strings(lines)
			}
			for _, line := range lines {
				fmt.Println(line)
			}
			return nil
		},
	}
}
