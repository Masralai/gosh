package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/urfave/cli/v3"
)

func Grep() *cli.Command {
	return &cli.Command{
		Name:      "grep",
		Usage:     "Search Text Using Patterns",
		UsageText: "cli grep 'pattern' <filename>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "f",
				Usage: "Search ignoring case differences (uppercase or lowercase)",
			},
			&cli.BoolFlag{
				Name:  "r",
				Usage: "recursive grep",
			},
			&cli.BoolFlag{
				Name:  "v",
				Usage: "Find lines that do not match the pattern",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("usage: grep 'pattern' ")
			}

			filename := c.Args().Get(1)
			pattern := c.Args().Get(0)

			var regObj *regexp.Regexp
			var err error
			if c.Bool("f") {
				regObj, err = regexp.Compile("(?i)" + pattern)
			} else {
				regObj, err = regexp.Compile(pattern)
			}
			if err != nil {
				return fmt.Errorf("Failed to create regex Object: %v", err)
			}

			if c.Bool("r") {
				return filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if d.IsDir() {
						return nil
					}
					file, err := os.Open(path)
					if err != nil {
						return nil
					}
					defer file.Close()

					scanner := bufio.NewScanner(file)
					lineNum := 0
					for scanner.Scan() {
						lineNum++
						line := scanner.Text()
						matched := regObj.MatchString(line)
						if c.Bool("v") {
							matched = !matched
						}
						if matched {
							fmt.Printf("%s:%d:%s\n", path, lineNum, line)
						}
					}
					return scanner.Err()
				})
			}

			root, _ := os.OpenRoot(".")
			defer root.Close()
			file, err := root.Open(filename)
			if err != nil {
				return fmt.Errorf("Failed to open file %v", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				matched := regObj.MatchString(line)
				if c.Bool("v") {
					matched = !matched
				}
				if matched {
					fmt.Printf("%s\n", line)
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
		UsageText: "cli head [-n <lines>] ",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "n",
				Usage: "-n",
				Value: 10,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("usage: head ")
			}
			file, err := os.Open(c.Args().Get(0))
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			count := 0
			maxLines := c.Int("n")
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
		Usage:     "Display Last Part of Files",
		UsageText: "cli tail [-n <lines>] ",
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
				return fmt.Errorf("usage: tail ")
			}
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
			n := c.Int("lines")
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
