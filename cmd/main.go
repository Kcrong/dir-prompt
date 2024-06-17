package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const (
	delimiter = "-------------"
)

func main() {
	// Define flags
	regexFlag := flag.String("regex", ".*", "Regex string for filtering files")
	rootFlag := flag.String("root", ".", "Base path to start file listing")
	flag.Parse()

	// Compile the regex
	regex, err := regexp.Compile(*regexFlag)
	if err != nil {
		fmt.Println("Invalid regex:", err)
		return
	}

	// List files recursively and filter
	files := listFiles(*rootFlag, regex)

	// Format and print the output
	printFiles(files)
}

func listFiles(root string, regex *regexp.Regexp) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && regex.MatchString(info.Name()) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error listing files:", err)
	}

	return files
}

func printFiles(files []string) {
	for i, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		if i > 0 {
			fmt.Println(delimiter)
		}
		fmt.Printf("Filename: %s\n%s\n", file, string(content[:]))
	}
}
