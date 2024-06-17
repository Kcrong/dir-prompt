# dirprompt
Go-based tool that scans the current directory, parses files, and outputs their filenames and contents in a structured format suitable for Large Language Model (LLM) prompts. Making it easy to integrate into various workflows that require formatted file content for AI processing.

## Features
**Directory Scanning**: Automatically scans the current directory for files.
**Content Parsing**: Reads and formats the contents of each file.
**LLM Prompt Formatting**: Outputs filenames and contents in a structured format ideal for LLM inputs.

## Usage
### Installation
```bash
go install github.com/Kcrong/dirprompt/cmd/dirprompt@main
```

### Run
```bash
dirprompt --regex ".*.go" --root "./"
```

### Output
```text
Filename: .github/workflows/go.yml
# This workflow will build a golang project
# For more information see: https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go

name: Go

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:

  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.22'

    - name: Build
      run: go build -o build/main -v cmd/dirprompt/main.go

    - name: Upload a Build Artifact
      uses: actions/upload-artifact@v4.3.3
      with:
        name: main
        path: build/main
-------------
Filename: .golangci.yml
run:
  timeout: 3m

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - zerologlint
    - gofmt
    - gocritic
    - revive
    - gci

linters-settings:
  gci:
    sections:
      - standard # Standard section: captures all standard packages.
      - default # Default section: contains all imports that could not be matched to another section type.
      - prefix(github.com/Kcrong/dirprompt) # Custom section: groups all imports with the specified Prefix.
      - blank # Blank section: contains all blank imports. This section is not present unless explicitly enabled.
      - dot # Dot section: contains all dot imports. This section is not present unless explicitly enabled.
      - alias # Alias section: contains all alias imports. This section is not present unless explicitly enabled.
    # Skip generated files.
    skip-generated: true

-------------
Filename: cmd/dirprompt/main.go
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

-------------
Filename: go.mod
module github.com/Kcrong/dirprompt

go 1.22.4
```
