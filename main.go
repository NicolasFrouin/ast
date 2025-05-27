package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

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

// Token represents a lexical token
type Token struct {
	Type  string
	Value string
}

// Lexer breaks input into tokens
type Lexer struct {
	input  string
	pos    int
	tokens []Token
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		tokens: []Token{},
	}
}

// Tokenize generates all tokens from the input
func (l *Lexer) Tokenize() []Token {
	for l.pos < len(l.input) {
		char := l.input[l.pos]

		// Skip whitespace
		if unicode.IsSpace(rune(char)) {
			l.pos++
			continue
		}

		// Handle comments
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
		if isAlpha(char) {
			start := l.pos
			for l.pos < len(l.input) && (isAlpha(l.input[l.pos]) || isDigit(l.input[l.pos])) {
				l.pos++
			}
			value := l.input[start:l.pos]

			// Check if it's a keyword
			tokenType := "IDENTIFIER"
			switch value {
			case "function":
				tokenType = "FUNCTION"
			case "return":
				tokenType = "RETURN"
			case "const":
				tokenType = "CONST"
			case "let":
				tokenType = "LET"
			case "var":
				tokenType = "VAR"
			}

			l.tokens = append(l.tokens, Token{Type: tokenType, Value: value})
			continue
		}

		// Handle strings
		if char == '"' || char == '\'' {
			quote := char
			start := l.pos
			l.pos++ // Skip the opening quote
			for l.pos < len(l.input) && l.input[l.pos] != quote {
				l.pos++
			}
			l.pos++ // Skip the closing quote
			l.tokens = append(l.tokens, Token{Type: "STRING", Value: l.input[start:l.pos]})
			continue
		}

		// Handle special characters
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
		case '=':
			l.tokens = append(l.tokens, Token{Type: "EQUALS", Value: "="})
		default:
			// Skip unknown characters
			l.pos++
			continue
		}
		l.pos++
	}

	l.tokens = append(l.tokens, Token{Type: "EOF", Value: ""})
	return l.tokens
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// AST node types
type Node interface {
	Type() string
}

type Program struct {
	Body []Node
}

func (p *Program) Type() string {
	return "Program"
}

type FunctionDeclaration struct {
	Name   string
	Params []string
	Body   []Node
}

func (f *FunctionDeclaration) Type() string {
	return "FunctionDeclaration"
}

type ReturnStatement struct {
	Argument Node
}

func (r *ReturnStatement) Type() string {
	return "ReturnStatement"
}

type Identifier struct {
	Name string
}

func (i *Identifier) Type() string {
	return "Identifier"
}

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Type() string {
	return "StringLiteral"
}

type VariableDeclaration struct {
	Kind  string // const, let, var
	Name  string
	Value Node
}

func (v *VariableDeclaration) Type() string {
	return "VariableDeclaration"
}

type Comment struct {
	Text string
}

func (c *Comment) Type() string {
	return "Comment"
}

// Parser generates an AST from tokens
type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: "EOF", Value: ""}
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() Token {
	p.pos++
	return p.current()
}

func (p *Parser) Parse() *Program {
	program := &Program{Body: []Node{}}

	for p.current().Type != "EOF" {
		node := p.parseStatement()
		if node != nil {
			program.Body = append(program.Body, node)
		}
	}

	return program
}

func (p *Parser) parseStatement() Node {
	token := p.current()

	switch token.Type {
	case "COMMENT":
		return p.parseComment()
	case "FUNCTION":
		return p.parseFunctionDeclaration()
	case "RETURN":
		return p.parseReturnStatement()
	case "CONST", "LET", "VAR":
		return p.parseVariableDeclaration()
	case "SEMICOLON":
		p.next()
		return nil
	default:
		// Unexpected token, skip it
		p.next()
		return nil
	}
}

func (p *Parser) parseComment() *Comment {
	comment := &Comment{Text: p.current().Value}
	p.next() // Skip comment token
	return comment
}

func (p *Parser) parseFunctionDeclaration() *FunctionDeclaration {
	p.next() // Skip function keyword

	name := p.current().Value
	p.next() // Skip identifier

	// Parse parameters
	params := []string{}
	if p.current().Type == "LEFT_PAREN" {
		p.next() // Skip (
		for p.current().Type != "RIGHT_PAREN" {
			if p.current().Type == "IDENTIFIER" {
				params = append(params, p.current().Value)
			}
			p.next()
		}
		p.next() // Skip )
	}

	// Parse function body
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

func (p *Parser) parseReturnStatement() *ReturnStatement {
	p.next() // Skip return keyword

	var argument Node
	if p.current().Type == "IDENTIFIER" {
		argument = &Identifier{Name: p.current().Value}
		p.next()
	}

	// Skip semicolon if present
	if p.current().Type == "SEMICOLON" {
		p.next()
	}

	return &ReturnStatement{Argument: argument}
}

func (p *Parser) parseVariableDeclaration() *VariableDeclaration {
	kind := p.current().Value
	p.next() // Skip const/let/var

	name := p.current().Value
	p.next() // Skip identifier

	// Skip equals sign
	if p.current().Type == "EQUALS" {
		p.next()
	}

	var value Node
	if p.current().Type == "STRING" {
		// Remove quotes from string literal
		rawValue := p.current().Value
		cleanValue := strings.Trim(rawValue, "\"'")
		value = &StringLiteral{Value: cleanValue}
		p.next()
	}

	// Skip semicolon if present
	if p.current().Type == "SEMICOLON" {
		p.next()
	}

	return &VariableDeclaration{Kind: kind, Name: name, Value: value}
}

// PrintAST pretty prints the AST
func PrintAST(node Node, indent string) {
	switch n := node.(type) {
	case *Program:
		fmt.Println(indent + "Program:")
		for _, stmt := range n.Body {
			PrintAST(stmt, indent+"  ")
		}
	case *FunctionDeclaration:
		fmt.Printf("%sFunctionDeclaration: %s\n", indent, n.Name)
		fmt.Printf("%s  Parameters: %v\n", indent, n.Params)
		fmt.Printf("%s  Body:\n", indent)
		for _, stmt := range n.Body {
			PrintAST(stmt, indent+"    ")
		}
	case *ReturnStatement:
		fmt.Printf("%sReturnStatement:\n", indent)
		if n.Argument != nil {
			PrintAST(n.Argument, indent+"  ")
		}
	case *Identifier:
		fmt.Printf("%sIdentifier: %s\n", indent, n.Name)
	case *StringLiteral:
		fmt.Printf("%sStringLiteral: %s\n", indent, n.Value)
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

func main() {
	content := readFile("./script.js")

	fmt.Println("File content:")
	fmt.Println(content)
	fmt.Println("\nTokenizing...")

	lexer := NewLexer(content)
	tokens := lexer.Tokenize()

	fmt.Println("\nTokens:")
	for _, token := range tokens {
		if token.Type != "EOF" {
			fmt.Printf("  %s: %s\n", token.Type, token.Value)
		}
	}

	fmt.Println("\nParsing...")
	parser := NewParser(tokens)
	ast := parser.Parse()

	fmt.Println("\nAST:")
	PrintAST(ast, "")
}
