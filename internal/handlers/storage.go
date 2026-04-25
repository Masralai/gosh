package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Zip() *cli.Command {
	return &cli.Command{

		Name:      "zip",
		Usage:     "compress file",
		UsageText: "cli zip <filename>",
		Action: func(ctx context.Context, c *cli.Command) error {
			archive, err := os.Create(c.Args().Get(0) + ".zip")
			if err != nil {
				return fmt.Errorf("failed to create zip archive :%v", err)
			}
			defer archive.Close()

			zipWriter := zip.NewWriter(archive)
			fmt.Println("opening first file...")
			f1, err := os.Open(c.Args().Get(1))
			if err != nil {
				return fmt.Errorf("file error:%v", err)
			}
			defer f1.Close()
			fmt.Println("adding file to the archive...")

			//compression path
			path := c.Args().Get(3)
			w1, err := zipWriter.Create(path)
			if err != nil {
				return fmt.Errorf("Failed to add file to archive:%v", err)
			}

			//copy uncompressed file to archive
			if _, err := io.Copy(w1, f1); err != nil {
				return fmt.Errorf("Failed to copy uncompressed file to archive:%v", err)
			}

			fmt.Println("opening second file ...")
			f2, err := os.Open(c.Args().Get(2))
			if err != nil {
				return fmt.Errorf("failed to open second file:%v", err)
			}
			defer f2.Close()

			//create entry in the zip archive
			w2, err := zipWriter.Create(path)
			if err != nil {
				return fmt.Errorf("Failed to add file to archive:%v", err)
			}

			//copy uncompressed file to archive
			if _, err := io.Copy(w2, f2); err != nil {
				return fmt.Errorf("Failed to copy uncompressed file to archive:%v", err)
			}

			if err := zipWriter.Close(); err != nil {
				return fmt.Errorf("failed to close zipwriter:%v", err)
			}
			return nil
		},
	}
}

func Unzip() *cli.Command {
	return &cli.Command{
		Name:      "unzip",
		Usage:     "Extract from ZIP archive",
		UsageText: "cli unzip <filename>.zip dest",
		Action: func(ctx context.Context, c *cli.Command) error {
			fmt.Println("opening zip archive")
			filename := c.Args().Get(0)
			archive, err := zip.OpenReader(filename)
			if err != nil {
				return fmt.Errorf("failed to read archive: %v", err)
			}
			defer archive.Close()

			dest := c.Args().Get(1)
			root, err := os.OpenRoot(dest)
			if err != nil {
				return fmt.Errorf("failed to open dest root:%v", err)
			}
			defer root.Close()

			for _, f := range archive.File {
				// filePath := filepath.Join(dest, f.Name)
				filePath := f.Name

				destAbs, _ := filepath.Abs(dest)
				fileAbs, _ := filepath.Abs(filePath)
				if !strings.HasPrefix(fileAbs, destAbs) {
					return fmt.Errorf("illegal file path:%s", f.Name)
				}

				fmt.Println("unzipping file...", filePath)
				//Decompression Bomb prevention
				if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
					return fmt.Errorf("invalid file path %s", filePath)
				}

				//empty dir
				if f.FileInfo().IsDir() {
					fmt.Println("creating directory")
					//os.ModePerm to 0750(User=All, Group=Read/Execute, Others=None).
					if err := root.MkdirAll(filePath, 0750); err != nil {
						return fmt.Errorf("failed to crreate empty dir: %v", err)
					}
					continue
				}

				//file within dir
				if err := root.MkdirAll(filepath.Dir(filePath), 0750); err != nil {
					return fmt.Errorf("failed to unzip :%v", err)
				}

				//read-write, create, trucate config
				destFile, err := root.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					return fmt.Errorf("failed to create empty dest: %v", err)
				}
				defer destFile.Close()

				//open file and copy contents
				fileInArchive, err := f.Open()
				if err != nil {
					return fmt.Errorf("failed to open file:%v", err)
				}
				defer fileInArchive.Close()

				const MxDecompress = 500 * 1024 * 1024 //500MiB

				if _, err := io.CopyN(destFile, fileInArchive, MxDecompress); err != nil && err != io.EOF {
					return fmt.Errorf("failed to copy contents or file too large: %v", err)
				}
			}
			return nil
		},
	}
}
