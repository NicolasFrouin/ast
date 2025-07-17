package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// main is the entry point of our program
// It reads the JS file, tokenizes it, builds the AST, and prints the result
func main() {
	// Define command-line flags
	filePath := flag.String("f", "./script.js", "Path to JavaScript file to parse")

	// Parse the command-line flags
	flag.Parse()

	// Validate the file extension
	if !strings.HasSuffix(*filePath, ".js") {
		fmt.Fprintf(os.Stderr, "Error: File must be a JavaScript file with .js extension\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -f <file.js>\n", os.Args[0])
		os.Exit(1)
	}

	// Check if the file exists
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File %s does not exist\n", *filePath)
		os.Exit(1)
	}

	// Read the JavaScript file
	content := readFile(*filePath)

	// Print file information
	fmt.Printf("Parsing file: %s\n", *filePath)
	fmt.Println("File content:")
	fmt.Println(content)
	fmt.Println("\nTokenizing...")

	// Tokenize the source code
	lexer := NewLexer(content)
	tokens := lexer.Tokenize()

	// Print all identified tokens for debugging
	fmt.Println("\nTokens:")
	for _, token := range tokens {
		if token.Type != "EOF" {
			fmt.Printf("  %s: %s\n", token.Type, token.Value)
		}
	}

	// Parse the tokens into an AST
	fmt.Println("\nParsing...")
	parser := NewParser(tokens)
	ast := parser.Parse()

	// Print the structure of the AST
	fmt.Println("\nAST:")
	PrintAST(ast, "")
}
