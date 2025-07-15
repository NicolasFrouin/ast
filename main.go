package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
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

// Token represents a lexical token in our JavaScript parser
// Type is the token category (like "FUNCTION", "IDENTIFIER", etc.)
// Value stores the actual text from the source code
type Token struct {
	Type  string
	Value string
}

// Lexer breaks input source code into tokens
// It scans through the input character by character to identify tokens
type Lexer struct {
	input  string  // The full source code text being analyzed
	pos    int     // Current position in the input (points to current character)
	tokens []Token // Collection of tokens found so far
}

// NewLexer creates a new lexer instance with the given input
// This is a constructor function that initializes a lexer ready for tokenization
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,         // Start at the beginning of input
		tokens: []Token{}, // Empty token list
	}
}

// Tokenize processes the entire input and converts it to tokens
// This is the main lexical analysis function that identifies all tokens in the source
func (l *Lexer) Tokenize() []Token {
	// Loop through the entire input
	for l.pos < len(l.input) {
		char := l.input[l.pos]

		// Skip whitespace (spaces, tabs, newlines)
		// Whitespace generally has no semantic meaning in JavaScript
		if unicode.IsSpace(rune(char)) {
			l.pos++
			continue
		}

		// Handle single-line comments (// comment)
		// Comments are preserved in our AST for documentation purposes
		if char == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
			// Skip to end of line
			start := l.pos
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			l.tokens = append(l.tokens, Token{Type: "COMMENT", Value: l.input[start:l.pos]})
			continue
		}

		// Handle identifiers and keywords
		// Identifiers include variable names, function names, etc.
		// Keywords are reserved words like 'function', 'return', etc.
		if isAlpha(char) {
			start := l.pos
			// Collect all alphanumeric characters that form this identifier
			for l.pos < len(l.input) && (isAlpha(l.input[l.pos]) || isDigit(l.input[l.pos])) {
				l.pos++
			}
			value := l.input[start:l.pos]

			// Check if the identifier is actually a keyword
			tokenType := "IDENTIFIER"
			switch value {
			case "function":
				tokenType = "FUNCTION" // Function declaration keyword
			case "return":
				tokenType = "RETURN" // Return statement keyword
			case "const":
				tokenType = "CONST" // Constant variable declaration
			case "let":
				tokenType = "LET" // Block-scoped variable declaration
			case "var":
				tokenType = "VAR" // Function-scoped variable declaration
			case "if":
				tokenType = "IF" // If statement keyword
			}

			l.tokens = append(l.tokens, Token{Type: tokenType, Value: value})
			continue
		}

		// Handle string literals ("string" or 'string')
		if char == '"' || char == '\'' {
			quote := char
			start := l.pos
			l.pos++ // Skip the opening quote
			// Continue until finding the matching closing quote
			for l.pos < len(l.input) && l.input[l.pos] != quote {
				l.pos++
			}
			l.pos++ // Skip the closing quote
			l.tokens = append(l.tokens, Token{Type: "STRING", Value: l.input[start:l.pos]})
			continue
		}

		// Handle special characters and syntax elements
		switch char {
		case '(':
			l.tokens = append(l.tokens, Token{Type: "LEFT_PAREN", Value: "("})
		case ')':
			l.tokens = append(l.tokens, Token{Type: "RIGHT_PAREN", Value: ")"})
		case '{':
			l.tokens = append(l.tokens, Token{Type: "LEFT_BRACE", Value: "{"})
		case '}':
			l.tokens = append(l.tokens, Token{Type: "RIGHT_BRACE", Value: "}"})
		case ';':
			l.tokens = append(l.tokens, Token{Type: "SEMICOLON", Value: ";"})
		case ',':
			l.tokens = append(l.tokens, Token{Type: "COMMA", Value: ","})
		case '=':
			// Check for equality operator (==)
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
				l.tokens = append(l.tokens, Token{Type: "EQUALITY", Value: "=="})
				l.pos++ // Skip the next '=' since we're handling both at once
			} else {
				l.tokens = append(l.tokens, Token{Type: "EQUALS", Value: "="})
			}
		case '>':
			// Check for greater than or equal (>=)
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
				l.tokens = append(l.tokens, Token{Type: "GREATER_EQUAL", Value: ">="})
				l.pos++ // Skip the next '=' since we're handling both at once
			} else {
				l.tokens = append(l.tokens, Token{Type: "GREATER_THAN", Value: ">"})
			}
		case '<':
			// Check for less than or equal (<=)
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
				l.tokens = append(l.tokens, Token{Type: "LESS_EQUAL", Value: "<="})
				l.pos++ // Skip the next '=' since we're handling both at once
			} else {
				l.tokens = append(l.tokens, Token{Type: "LESS_THAN", Value: "<"})
			}
		case '+':
			l.tokens = append(l.tokens, Token{Type: "PLUS", Value: "+"})
		case '-':
			l.tokens = append(l.tokens, Token{Type: "MINUS", Value: "-"})
		case '*':
			l.tokens = append(l.tokens, Token{Type: "MULTIPLY", Value: "*"})
		case '/':
			// Check if it's a comment (already handled above) or division
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
				// This is a comment, skip it (already handled in comment section)
				l.pos++
				continue
			} else {
				l.tokens = append(l.tokens, Token{Type: "DIVIDE", Value: "/"})
			}
		case '%':
			l.tokens = append(l.tokens, Token{Type: "MODULO", Value: "%"})
		default:
			// Handle numeric literals (including decimals)
			if isDigit(char) {
				start := l.pos
				for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
					l.pos++
				}
				// Handle decimal part if present
				if l.pos < len(l.input) && l.input[l.pos] == '.' {
					l.pos++ // Skip the decimal point
					// Consume digits after decimal point
					for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
						l.pos++
					}
				}
				l.tokens = append(l.tokens, Token{Type: "NUMBER", Value: l.input[start:l.pos]})
				continue
			}

			// Skip unknown characters
			l.pos++
			continue
		}
		l.pos++
	}

	// Add an EOF (End Of File) token to indicate the end of input
	l.tokens = append(l.tokens, Token{Type: "EOF", Value: ""})
	return l.tokens
}

// isAlpha checks if a character is alphabetic or underscore
// Used to determine the start of identifiers
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isDigit checks if a character is a numeric digit
// Used for the non-first characters of identifiers
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// Node is an interface representing any node in our Abstract Syntax Tree
// Every AST node type must implement the Type method
type Node interface {
	Type() string
}

// Program is the root node of our Abstract Syntax Tree
// It contains all the top-level statements in the source file
type Program struct {
	Body []Node // Array of top-level statements
}

func (p *Program) Type() string {
	return "Program"
}

// FunctionDeclaration represents a JavaScript function definition
// Example: function name(param1, param2 = defaultValue) { ... }
type FunctionDeclaration struct {
	Name   string      // Function name
	Params []Parameter // Parameter names and default values
	Body   []Node      // Function body statements
}

func (f *FunctionDeclaration) Type() string {
	return "FunctionDeclaration"
}

// Parameter represents a function parameter with optional default value
type Parameter struct {
	Name         string // Parameter name
	DefaultValue Node   // Default value (nil if no default)
}

// ReturnStatement represents a 'return' statement in JavaScript
// Example: return expression;
type ReturnStatement struct {
	Argument Node // The value being returned (can be nil)
}

func (r *ReturnStatement) Type() string {
	return "ReturnStatement"
}

// Identifier represents a variable or function name
// Examples: x, myFunction, etc.
type Identifier struct {
	Name string // The name of the identifier
}

func (i *Identifier) Type() string {
	return "Identifier"
}

// StringLiteral represents a string value in the code
// Examples: "hello", 'world'
type StringLiteral struct {
	Value string // The actual string value without quotes
}

func (s *StringLiteral) Type() string {
	return "StringLiteral"
}

// VariableDeclaration represents a variable declaration
// Examples: const x = 5; let name = "value";
type VariableDeclaration struct {
	Kind  string // Declaration type: "const", "let", or "var"
	Name  string // Variable name
	Value Node   // Initial value (can be nil)
}

func (v *VariableDeclaration) Type() string {
	return "VariableDeclaration"
}

// Comment represents a code comment
// Example: // This is a comment
type Comment struct {
	Text string // The full text of the comment including //
}

func (c *Comment) Type() string {
	return "Comment"
}

// IfStatement represents an if conditional statement
// Example: if (condition) { ... }
type IfStatement struct {
	Test       Node   // The condition being tested
	Consequent []Node // Statements to execute if condition is true
}

func (i *IfStatement) Type() string {
	return "IfStatement"
}

// BinaryExpression represents expressions with two operands and an operator
// Examples: a == b, x + y
type BinaryExpression struct {
	Left     Node   // Left operand
	Operator string // Operator (e.g., "==", "+")
	Right    Node   // Right operand
}

func (b *BinaryExpression) Type() string {
	return "BinaryExpression"
}

// NumericLiteral represents numeric values in the code
// Example: 1, 3.14
type NumericLiteral struct {
	Value string // The numeric value
}

func (n *NumericLiteral) Type() string {
	return "NumericLiteral"
}

// Parser generates an AST from tokens
// It implements a recursive descent parser pattern
type Parser struct {
	tokens []Token // Token stream from the lexer
	pos    int     // Current position in the token stream
}

// NewParser creates a new parser with the given token stream
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// current returns the current token without advancing
func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: "EOF", Value: ""} // Return EOF if we're past the end
	}
	return p.tokens[p.pos]
}

// next moves to the next token and returns it
func (p *Parser) next() Token {
	p.pos++
	return p.current()
}

// Parse builds a complete AST from the token stream
// This is the entry point to the parsing process
func (p *Parser) Parse() *Program {
	program := &Program{Body: []Node{}}

	// Process tokens until EOF
	for p.current().Type != "EOF" {
		node := p.parseStatement()
		if node != nil {
			program.Body = append(program.Body, node)
		}
	}

	return program
}

// parseStatement parses a single statement based on the current token
// Different token types lead to different statement types
func (p *Parser) parseStatement() Node {
	token := p.current()

	switch token.Type {
	case "COMMENT":
		return p.parseComment() // Handle comments
	case "FUNCTION":
		return p.parseFunctionDeclaration() // Handle function declarations
	case "RETURN":
		return p.parseReturnStatement() // Handle return statements
	case "CONST", "LET", "VAR":
		return p.parseVariableDeclaration() // Handle variable declarations
	case "IF":
		return p.parseIfStatement() // Handle if statements
	case "SEMICOLON":
		p.next() // Skip standalone semicolons
		return nil
	default:
		// Skip tokens we don't recognize in this context
		p.next()
		return nil
	}
}

// parseComment creates a Comment node from a comment token
func (p *Parser) parseComment() *Comment {
	comment := &Comment{Text: p.current().Value}
	p.next() // Skip comment token
	return comment
}

// parseFunctionDeclaration parses a function declaration statement
// Format: function name(param1, param2 = defaultValue) { body }
func (p *Parser) parseFunctionDeclaration() *FunctionDeclaration {
	p.next() // Skip function keyword

	name := p.current().Value
	p.next() // Skip identifier

	// Parse parameters inside parentheses
	params := []Parameter{}
	if p.current().Type == "LEFT_PAREN" {
		p.next() // Skip (
		for p.current().Type != "RIGHT_PAREN" && p.current().Type != "EOF" {
			if p.current().Type == "IDENTIFIER" {
				paramName := p.current().Value
				p.next() // Skip parameter name

				var defaultValue Node
				// Check for default value assignment
				if p.current().Type == "EQUALS" {
					p.next() // Skip the equals sign
					defaultValue = p.parseExpression()
				}

				params = append(params, Parameter{
					Name:         paramName,
					DefaultValue: defaultValue,
				})

				// Skip comma if present
				if p.current().Type == "COMMA" {
					p.next()
				}
			} else {
				p.next() // Skip unexpected tokens
			}
		}
		if p.current().Type == "RIGHT_PAREN" {
			p.next() // Skip )
		}
	}

	// Parse function body inside braces
	body := []Node{}
	if p.current().Type == "LEFT_BRACE" {
		p.next() // Skip {
		for p.current().Type != "RIGHT_BRACE" {
			stmt := p.parseStatement()
			if stmt != nil {
				body = append(body, stmt)
			}
		}
		p.next() // Skip }
	}

	return &FunctionDeclaration{Name: name, Params: params, Body: body}
}

// parseIfStatement parses an if statement
// Format: if (condition) { body }
func (p *Parser) parseIfStatement() *IfStatement {
	p.next() // Skip the 'if' keyword

	// Parse condition in parentheses
	var test Node
	if p.current().Type == "LEFT_PAREN" {
		p.next() // Skip the opening parenthesis
		test = p.parseExpression()

		// Skip the closing parenthesis if present
		if p.current().Type == "RIGHT_PAREN" {
			p.next()
		}
	}

	// Parse consequent (the "then" block)
	consequent := []Node{}
	if p.current().Type == "LEFT_BRACE" {
		p.next() // Skip the opening brace
		// Parse statements until we reach the closing brace
		for p.current().Type != "RIGHT_BRACE" && p.current().Type != "EOF" {
			stmt := p.parseStatement()
			if stmt != nil {
				consequent = append(consequent, stmt)
			}
		}
		if p.current().Type == "RIGHT_BRACE" {
			p.next() // Skip the closing brace
		}
	}

	return &IfStatement{
		Test:       test,
		Consequent: consequent,
	}
}

// parseExpression parses expressions like comparisons and math operations
func (p *Parser) parseExpression() Node {
	// Parse the left side of the expression
	left := p.parsePrimary()

	// If followed by an operator, it's a binary expression
	if isBinaryOperator(p.current().Type) {
		operator := p.current().Value
		p.next() // Skip the operator

		// Parse the right side of the expression
		right := p.parsePrimary()

		return &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

// parsePrimary parses a primary expression (identifiers, literals)
func (p *Parser) parsePrimary() Node {
	token := p.current()

	switch token.Type {
	case "IDENTIFIER":
		identifier := &Identifier{Name: token.Value}
		p.next()
		return identifier
	case "NUMBER":
		number := &NumericLiteral{Value: token.Value}
		p.next()
		return number
	case "STRING":
		// Remove quotes from string literal
		rawValue := token.Value
		cleanValue := strings.Trim(rawValue, "\"'")
		value := &StringLiteral{Value: cleanValue}
		p.next()
		return value
	default:
		p.next() // Skip unhandled tokens
		return nil
	}
}

// isBinaryOperator checks if a token type represents a binary operator
func isBinaryOperator(tokenType string) bool {
	return tokenType == "EQUALITY" || tokenType == "EQUALS" ||
		tokenType == "PLUS" || tokenType == "MINUS" ||
		tokenType == "MULTIPLY" || tokenType == "DIVIDE" ||
		tokenType == "MODULO" ||
		tokenType == "GREATER_THAN" || tokenType == "LESS_THAN" ||
		tokenType == "GREATER_EQUAL" || tokenType == "LESS_EQUAL"
}

// parseReturnStatement parses a return statement
// Format: return expression;
func (p *Parser) parseReturnStatement() *ReturnStatement {
	p.next() // Skip return keyword

	var argument Node
	// Parse any expression as the return value
	// This handles: identifiers, literals, binary expressions, etc.
	if p.current().Type != "SEMICOLON" && p.current().Type != "EOF" {
		argument = p.parseExpression()
	}

	// Skip semicolon if present
	if p.current().Type == "SEMICOLON" {
		p.next()
	}

	return &ReturnStatement{Argument: argument}
}

// parseVariableDeclaration parses a variable declaration
// Format: const/let/var name = value;
func (p *Parser) parseVariableDeclaration() *VariableDeclaration {
	kind := p.current().Value
	p.next() // Skip const/let/var

	name := p.current().Value
	p.next() // Skip identifier

	// Skip equals sign
	if p.current().Type == "EQUALS" {
		p.next()
	}

	// Always try to parse as an expression first, which handles all cases:
	// - Simple literals (strings, numbers, identifiers)
	// - Complex expressions (1 + 2, a * b, etc.)
	value := p.parseExpression()

	// Skip semicolon if present
	if p.current().Type == "SEMICOLON" {
		p.next()
	}

	return &VariableDeclaration{Kind: kind, Name: name, Value: value}
}

// PrintAST recursively prints the AST in a human-readable format
// It uses indentation to show the tree structure
func PrintAST(node Node, indent string) {
	switch n := node.(type) {
	case *Program:
		fmt.Println(indent + "Program:")
		for _, stmt := range n.Body {
			PrintAST(stmt, indent+"  ")
		}
	case *FunctionDeclaration:
		fmt.Printf("%sFunctionDeclaration: %s\n", indent, n.Name)
		fmt.Printf("%s  Parameters:\n", indent)
		for _, param := range n.Params {
			if param.DefaultValue != nil {
				fmt.Printf("%s    %s (default):\n", indent, param.Name)
				PrintAST(param.DefaultValue, indent+"      ")
			} else {
				fmt.Printf("%s    %s\n", indent, param.Name)
			}
		}
		fmt.Printf("%s  Body:\n", indent)
		for _, stmt := range n.Body {
			PrintAST(stmt, indent+"    ")
		}
	case *IfStatement:
		fmt.Printf("%sIfStatement:\n", indent)
		fmt.Printf("%s  Condition:\n", indent)
		PrintAST(n.Test, indent+"    ")
		fmt.Printf("%s  Body:\n", indent)
		for _, stmt := range n.Consequent {
			PrintAST(stmt, indent+"    ")
		}
	case *BinaryExpression:
		fmt.Printf("%sBinaryExpression: %s\n", indent, n.Operator)
		fmt.Printf("%s  Left:\n", indent)
		PrintAST(n.Left, indent+"    ")
		fmt.Printf("%s  Right:\n", indent)
		PrintAST(n.Right, indent+"    ")
	case *ReturnStatement:
		fmt.Printf("%sReturnStatement:\n", indent)
		if n.Argument != nil {
			PrintAST(n.Argument, indent+"  ")
		}
	case *Identifier:
		fmt.Printf("%sIdentifier: %s\n", indent, n.Name)
	case *StringLiteral:
		fmt.Printf("%sStringLiteral: %s\n", indent, n.Value)
	case *NumericLiteral:
		fmt.Printf("%sNumericLiteral: %s\n", indent, n.Value)
	case *VariableDeclaration:
		fmt.Printf("%sVariableDeclaration: %s %s\n", indent, n.Kind, n.Name)
		if n.Value != nil {
			PrintAST(n.Value, indent+"  ")
		}
	case *Comment:
		fmt.Printf("%sComment: %s\n", indent, n.Text)
	default:
		fmt.Printf("%sUnknown node type\n", indent)
	}
}

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
