package main

import (
	"strings"
)

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
