package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func Zip() *cli.Command {
	return &cli.Command{
		Name:      "zip",
		Usage:     "Create a ZIP archive from files",
		UsageText: "zip <archive> <file1> [file2] ...",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("usage: zip <archive> <file1> [file2] ...")
			}

			archiveName := c.Args().Get(0)
			if !strings.HasSuffix(strings.ToLower(archiveName), ".zip") {
				archiveName += ".zip"
			}

			// #nosec G304 -- archive path is user-specified
			archive, err := os.Create(archiveName)
			if err != nil {
				return fmt.Errorf("failed to create zip archive: %v", err)
			}
			defer archive.Close()

			zipWriter := zip.NewWriter(archive)
			defer func() { _ = zipWriter.Close() }()

			for _, filePath := range c.Args().Slice()[1:] {
				cleanPath := filepath.Clean(filePath)
				info, err := os.Stat(cleanPath)
				if err != nil {
					return fmt.Errorf("failed to access %s: %v", filePath, err)
				}
				if info.IsDir() {
					return fmt.Errorf("zip of directories is not supported: %s", filePath)
				}
				// #nosec G304
				f, err := os.Open(cleanPath)
				if err != nil {
					return fmt.Errorf("failed to open %s: %v", filePath, err)
				}

				nameInZip := filepath.Base(cleanPath)
				w, err := zipWriter.Create(nameInZip)
				if err != nil {
					_ = f.Close()
					return fmt.Errorf("failed to add file to archive: %v", err)
				}

				if _, err := io.Copy(w, f); err != nil {
					_ = f.Close()
					return fmt.Errorf("failed to write %s to archive: %v", filePath, err)
				}
				if err := f.Close(); err != nil {
					return fmt.Errorf("failed to close %s: %v", filePath, err)
				}
			}

			return nil
		},
	}
}

func Unzip() *cli.Command {
	return &cli.Command{
		Name:      "unzip",
		Usage:     "Extract files from a ZIP archive",
		UsageText: "unzip <archive.zip> <destination>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("usage: unzip <archive.zip> <destination>")
			}
			filename := c.Args().Get(0)
			dest := c.Args().Get(1)

			archive, err := zip.OpenReader(filename)
			if err != nil {
				return fmt.Errorf("failed to open archive: %v", err)
			}
			defer archive.Close()

			if err := os.MkdirAll(dest, 0o750); err != nil {
				return fmt.Errorf("failed to create destination: %v", err)
			}
			root, err := os.OpenRoot(dest)
			if err != nil {
				return fmt.Errorf("failed to open destination: %v", err)
			}
			defer root.Close()

			for _, f := range archive.File {
				name := filepath.Clean(f.Name)
				name = filepath.ToSlash(name)
				if name == "." || name == "" {
					continue
				}
				if strings.HasPrefix(name, "/") || strings.HasPrefix(name, "../") || strings.Contains(name, "/../") || name == ".." {
					return fmt.Errorf("invalid file path in archive: %s", f.Name)
				}
				// normalize to OS separators for root ops; keep relative
				rel := filepath.FromSlash(strings.TrimPrefix(name, "./"))

				if f.FileInfo().IsDir() || strings.HasSuffix(f.Name, "/") {
					if err := root.MkdirAll(rel, 0o750); err != nil {
						return fmt.Errorf("failed to create directory %s: %v", f.Name, err)
					}
					continue
				}

				if dir := filepath.Dir(rel); dir != "." {
					if err := root.MkdirAll(dir, 0o750); err != nil {
						return fmt.Errorf("failed to create directory for %s: %v", f.Name, err)
					}
				}

				mode := f.Mode().Perm()
				if mode == 0 {
					mode = 0o644
				}
				destFile, err := root.OpenFile(rel, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
				if err != nil {
					return fmt.Errorf("failed to create %s: %v", f.Name, err)
				}

				srcFile, err := f.Open()
				if err != nil {
					_ = destFile.Close()
					return fmt.Errorf("failed to read %s from archive: %v", f.Name, err)
				}

				const maxDecompress = 500 * 1024 * 1024 // 500 MiB
				written, err := io.CopyN(destFile, srcFile, maxDecompress+1)
				_ = srcFile.Close()
				cerr := destFile.Close()
				if err != nil && err != io.EOF {
					return fmt.Errorf("failed to extract %s: %v", f.Name, err)
				}
				if written > maxDecompress {
					return fmt.Errorf("failed to extract %s: file too large", f.Name)
				}
				if cerr != nil {
					return fmt.Errorf("failed to write %s: %v", f.Name, cerr)
				}
			}
			return nil
		},
	}
}
