package main

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
