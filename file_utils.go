package main

import (
	"io"
	"os"
)

// readFile reads a file and returns its contents as a string
// It handles file opening and reading, with error handling
func readFile(filename string) string {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return string(content)
}
