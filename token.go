package main

import (
	"unicode"
)

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
