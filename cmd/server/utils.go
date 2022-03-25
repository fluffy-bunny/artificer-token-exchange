package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func DumpPath(root string) {
	files, err := FilePathWalkDir(root)
	if err != nil {
		panic(err)
	}
	fmt.Println("==========================================================")
	fmt.Println("Dumping files in " + root)
	fmt.Println("==========================================================")
	for _, file := range files {
		fmt.Println(file)
	}
	fmt.Println("==========================================================")
}
