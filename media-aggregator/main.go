package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	sourcePtr := flag.String("source", "", "Source directory (required)")
	destPtr := flag.String("dest", "", "Destination directory (required)")

	flag.Parse()

	if *sourcePtr == "" || *destPtr == "" {
		flag.Usage()
		return
	}

	sourceDir := *sourcePtr
	destDir := *destPtr

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isMediaFile(info.Name()) {
			err := copyFile(path, destDir)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Media files copied successfully!")
	}
}

func isMediaFile(filename string) bool {
	ext := filepath.Ext(filename)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif", ".mp4", ".mov", ".webm":
		return true
	default:
		return false
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(filepath.Join(dst, filepath.Base(src)))
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
