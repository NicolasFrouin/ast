# JavaScript AST Builder in Go

> A learning project: Building an Abstract Syntax Tree parser for JavaScript from scratch using Go

## Table of Contents

- [Overview](#overview)
- [Learning Journey](#learning-journey)
- [Core Concepts Discovered](#core-concepts-discovered)
- [Architecture](#architecture)
- [Implementation Details](#implementation-details)
- [Examples](#examples)
- [Usage](#usage)
- [Supported JavaScript Features](#supported-javascript-features)
- [What I Learned](#what-i-learned)

## Overview

This project is a JavaScript parser that builds an Abstract Syntax Tree (AST) from JavaScript source code, implemented entirely in Go. As a computer science student, I built this from scratch to understand how programming language parsers work under the hood.

The parser can handle basic JavaScript constructs including:

- Function declarations
- Variable declarations (`const`, `let`, `var`)
- If statements with conditions
- String and numeric literals
- Comments
- Binary expressions (equality comparisons)

## Learning Journey

### 1. Understanding What an AST Is

When I started this project, I first needed to understand what an Abstract Syntax Tree actually represents. An AST is a tree representation of the syntactic structure of source code, where:

- **Each node** represents a construct occurring in the programming language
- **The tree structure** shows the relationship between different parts of the code
- **Abstract** means it doesn't include every detail from the source (like whitespace or exact formatting)

For example, this JavaScript code:

```javascript
function greet(name) {
  return 'Hello ' + name;
}
```

Gets represented as a tree structure where the root is a Program node, containing a FunctionDeclaration node, which contains parameter information and a body with a ReturnStatement, etc.

### 2. The Two-Phase Approach

I discovered that parsing typically happens in two main phases:

#### Phase 1: Lexical Analysis (Tokenization)

Breaking the source code into **tokens** - the basic building blocks of the language.

#### Phase 2: Syntactic Analysis (Parsing)

Taking those tokens and building the actual tree structure.

This separation makes the code much more manageable and follows how real-world parsers work.

## Core Concepts Discovered

### Tokens - The Building Blocks

A token represents a meaningful sequence of characters in the source code. I learned that tokens have two main parts:

- **Type**: What kind of token it is (keyword, identifier, operator, etc.)
- **Value**: The actual text from the source code

```go
type Token struct {
    Type  string  // "FUNCTION", "IDENTIFIER", "STRING", etc.
    Value string  // "function", "myVar", "hello world", etc.
}
```

### The Lexer - Character-by-Character Analysis

The lexer is like a scanner that reads through the source code character by character and groups them into tokens. I implemented several key features:

1. **Whitespace handling**: Skip spaces, tabs, newlines
2. **Comment recognition**: Handle `//` style comments
3. **String literal parsing**: Handle both `"` and `'` quoted strings
4. **Keyword identification**: Recognize reserved words like `function`, `return`, `if`
5. **Operator recognition**: Handle `=`, `==`, parentheses, braces
6. **Number parsing**: Recognize numeric literals

### Parser - Building the Tree

The parser takes the stream of tokens and builds the actual tree structure. I learned about:

#### Recursive Descent Parsing

This is a top-down parsing technique where each grammar rule becomes a function. For example:

- `parseStatement()` - handles any kind of statement
- `parseFunctionDeclaration()` - specifically handles function declarations
- `parseExpression()` - handles expressions like comparisons

#### Grammar Rules

I had to define how different JavaScript constructs should be parsed:

```ebnf
FunctionDeclaration := "function" IDENTIFIER "(" ParameterList ")" "{" StatementList "}"
IfStatement := "if" "(" Expression ")" "{" StatementList "}"
VariableDeclaration := ("const"|"let"|"var") IDENTIFIER "=" Expression ";"
```

## Architecture

### 1. File Reading Layer

```go
func readFile(filename string) string
```

Simple file I/O to read JavaScript source files.

### 2. Lexical Analysis Layer

```go
type Lexer struct {
    input  string    // Source code
    pos    int       // Current position
    tokens []Token   // Generated tokens
}
```

The lexer processes the input character by character and generates a sequence of tokens.

### 3. AST Node Definitions

I created a type hierarchy using Go interfaces:

```go
type Node interface {
    Type() string
}
```

Then implemented specific node types for each JavaScript construct:

- `Program` - Root of the tree
- `FunctionDeclaration` - Function definitions
- `VariableDeclaration` - Variable declarations
- `IfStatement` - Conditional statements
- `ReturnStatement` - Return statements
- `BinaryExpression` - Operations like `==`
- `Identifier` - Variable/function names
- `StringLiteral` - String values
- `NumericLiteral` - Number values
- `Comment` - Code comments

### 4. Parser Layer

```go
type Parser struct {
    tokens []Token
    pos    int
}
```

The parser maintains the current position in the token stream and builds the AST using recursive descent.

### 5. Pretty Printing

A recursive function that traverses the AST and prints it in a human-readable format with proper indentation.

## Implementation Details

### Token Recognition Patterns

**Keywords**: I use a switch statement to identify reserved words:

```go
switch value {
case "function":
    tokenType = "FUNCTION"
case "return":
    tokenType = "RETURN"
// ... etc
}
```

**Identifiers**: Start with a letter or underscore, followed by letters, digits, or underscores:

```go
func isAlpha(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}
```

**String Literals**: Handle both single and double quotes:

```go
if char == '"' || char == '\'' {
    quote := char
    // ... collect until matching quote
}
```

### Parsing Strategy

I implemented a **recursive descent parser** where each grammar rule becomes a function:

1. **Top-level parsing**: `Parse()` repeatedly calls `parseStatement()`
2. **Statement dispatch**: `parseStatement()` looks at the current token type and calls the appropriate specific parser
3. **Specific parsers**: Each construct has its own parsing function that knows the expected token sequence

### Error Handling Philosophy

For this learning project, I chose simple error handling:

- File errors cause `panic()`
- Unknown tokens are skipped
- Malformed syntax results in `nil` nodes

In a production parser, I would implement proper error recovery and meaningful error messages.

## Examples

### Input JavaScript File (`script.js`)

```javascript
// comment

function funcName(funcArg) {
  if (funcArg == 1) {
    return 'Function argument is 1';
  }

  return funcArg;
}

const constVar = 'This is a constant variable';
```

### Generated Tokens

```text
COMMENT: // comment
FUNCTION: function
IDENTIFIER: funcName
LEFT_PAREN: (
IDENTIFIER: funcArg
RIGHT_PAREN: )
LEFT_BRACE: {
IF: if
LEFT_PAREN: (
IDENTIFIER: funcArg
EQUALITY: ==
NUMBER: 1
RIGHT_PAREN: )
LEFT_BRACE: {
RETURN: return
STRING: "Function argument is 1"
SEMICOLON: ;
RIGHT_BRACE: }
RETURN: return
IDENTIFIER: funcArg
SEMICOLON: ;
RIGHT_BRACE: }
CONST: const
IDENTIFIER: constVar
EQUALS: =
STRING: "This is a constant variable"
SEMICOLON: ;
EOF:
```

### Generated AST Structure

```text
Program:
  Comment: // comment
  FunctionDeclaration: funcName
    Parameters: [funcArg]
    Body:
      IfStatement:
        Condition:
          BinaryExpression: ==
            Left:
              Identifier: funcArg
            Right:
              NumericLiteral: 1
        Body:
          ReturnStatement:
            Argument:
              StringLiteral: Function argument is 1
      ReturnStatement:
        Argument:
          Identifier: funcArg
  VariableDeclaration: const constVar = "This is a constant variable"
```

## Usage

### Prerequisites

- Go 1.19 or higher
- A JavaScript file to parse

### Running the Parser

1. **Basic usage** (parses `./script.js` by default):

```bash
go run main.go
```

1. **Parse a specific file**:

```bash
go run main.go -f path/to/your/file.js
```

1. **Build and run**:

```bash
go build -o js-parser main.go
./js-parser -f script.js
```

### Command Line Options

- `-f <filepath>`: Specify the JavaScript file to parse (default: `./script.js`)

## Supported JavaScript Features

### ✅ Currently Supported

- **Function declarations**: `function name(params) { ... }`
- **Variable declarations**: `const`, `let`, `var` with string initialization
- **If statements**: `if (condition) { ... }` with equality comparisons
- **Return statements**: `return value;`
- **Comments**: `// single line comments`
- **String literals**: Both `"double"` and `'single'` quoted
- **Numeric literals**: Integer numbers
- **Binary expressions**: Equality operator `==`
- **Identifiers**: Variable and function names

### ❌ Not Yet Supported

- **Functions calls**: `myFunction()`
- **Arithmetic operators**: `+`, `-`, `*`, `/`
- **Else clauses**: `if ... else ...`
- **Loops**: `for`, `while`
- **Objects and arrays**: `{}`, `[]`
- **Arrow functions**: `() => {}`
- **Template literals**: `` `string ${var}` ``
- **Multiple variable declarations**: `let a, b, c;`

## What I Learned

### Technical Skills

1. **Lexical Analysis**: How to break text into meaningful tokens
2. **Parsing Theory**: Recursive descent parsing and grammar rules
3. **Go Language Mastery**: Deep dive into Go's type system and idioms (detailed below)
4. **Tree Data Structures**: Building and traversing hierarchical data
5. **Command-line Programs**: Using the `flag` package for CLI arguments

### Learning Go Through This Project

This project was an excellent vehicle for learning Go because it required me to use many of the language's core features in practical ways:

#### Interfaces and Polymorphism

One of the most important Go concepts I learned was how to use interfaces effectively. Coming from other languages, I initially struggled with Go's implicit interface satisfaction, but this project made it click:

```go
type Node interface {
    Type() string
}
```

Every AST node type implements this interface simply by having a `Type()` method. This allowed me to:

- Store different node types in the same slice: `[]Node`
- Write generic functions that work with any node type
- Use type assertions to access specific node properties when needed

The beauty is that I never had to explicitly declare that my structs implement the interface - they just do!

#### Structs and Methods

I learned how Go's struct system works differently from classes in other languages:

```go
type Lexer struct {
    input  string
    pos    int
    tokens []Token
}

func (l *Lexer) Tokenize() []Token {
    // Method with receiver
}
```

Key insights I gained:

- **Receiver types**: When to use pointer receivers (`*Lexer`) vs value receivers (`Lexer`)
- **Method organization**: How methods are associated with types outside the struct definition
- **Encapsulation**: Using capitalized names for public fields/methods, lowercase for private

#### Pointer vs Value Semantics

This project taught me when to use pointers in Go:

```go
// Constructor returns pointer because we'll modify the lexer
func NewLexer(input string) *Lexer {
    return &Lexer{
        input:  input,
        pos:    0,
        tokens: []Token{},
    }
}

// Method uses pointer receiver to modify the lexer's state
func (l *Lexer) Tokenize() []Token {
    // l.pos++ modifies the original lexer
}
```

I learned that:

- Large structs should typically be passed as pointers for performance
- Methods that modify state need pointer receivers
- Go's garbage collector handles pointer management automatically

#### Type Assertions and Type Switches

The AST printing function taught me about Go's type assertion system:

```go
func PrintAST(node Node, indent string) {
    switch n := node.(type) {
    case *Program:
        fmt.Println(indent + "Program:")
        for _, stmt := range n.Body {
            PrintAST(stmt, indent+"  ")
        }
    case *FunctionDeclaration:
        fmt.Printf("%sFunctionDeclaration: %s\n", indent, n.Name)
        // Access FunctionDeclaration-specific fields
    case *Identifier:
        fmt.Printf("%sIdentifier: %s\n", indent, n.Name)
    }
}
```

This taught me:

- How to safely convert interface types to concrete types
- The power of type switches for handling different implementations
- Go's approach to runtime type checking

#### Slices and Memory Management

Working with token streams taught me about Go's slice behavior:

```go
// Growing slices dynamically
l.tokens = append(l.tokens, Token{Type: "FUNCTION", Value: "function"})

// Slicing for substrings
value := l.input[start:l.pos]
```

I learned:

- How `append()` works and when it reallocates
- Slice capacity vs length
- Efficient string manipulation using slicing

#### Go's Error Handling Approach

While I chose simple error handling for this learning project, I discovered Go's explicit error handling approach:

```go
func readFile(filename string) string {
    file, err := os.Open(filename)
    if err != nil {
        panic(err) // Simple for learning, but not idiomatic
    }
    defer file.Close()
    
    content, err := io.ReadAll(file)
    if err != nil {
        panic(err)
    }
    
    return string(content)
}
```

This taught me:

- Go's explicit error handling (no exceptions)
- The importance of checking every error
- How `defer` works for cleanup
- Why Go developers prefer explicit error handling

#### Package Organization and Imports

I learned how Go's package system works:

```go
import (
    "flag"      // Standard library
    "fmt"       // Standard library  
    "io"        // Standard library
    "os"        // Standard library
    "strings"   // Standard library
    "unicode"   // Standard library
)
```

Key insights:

- How Go organizes code into packages
- The difference between local packages and standard library
- Import naming and organization conventions
- How `go mod` manages dependencies

#### Go's Approach to Object-Oriented Programming

This project showed me how Go does OOP differently:

- **No inheritance**: Use composition instead
- **No classes**: Use structs with methods
- **Interfaces**: Focus on behavior, not inheritance hierarchies
- **Embedding**: Achieve code reuse through struct embedding

#### Concurrency Considerations

While this parser is single-threaded, building it made me think about Go's concurrency model:

- How goroutines could parallelize lexing and parsing
- Where channels might be useful for token streaming
- Why Go's "share memory by communicating" philosophy matters

#### Performance Insights

Working character-by-character through source code taught me about Go's performance characteristics:

- String operations and garbage collection
- When to use `strings.Builder` vs string concatenation
- How Go's compiler optimizes certain patterns

#### Debugging and Development

I learned practical Go development skills:

- Using `fmt.Printf` for debugging
- How Go's compiler helps catch errors early
- The importance of `go fmt` for consistent formatting
- How to structure a Go project

This project was perfect for learning Go because it required:

- Complex data structures (trees, slices)
- Interface design and polymorphism
- String manipulation and character processing
- Method design and receiver types
- Package organization

By the end, I felt comfortable with Go's philosophy of simplicity and explicitness. The language's constraints actually made the parser easier to reason about - no hidden inheritance, explicit error handling, and clear data flow.

### Problem-Solving Approaches

1. **Divide and Conquer**: Breaking complex parsing into smaller, manageable pieces
2. **Incremental Development**: Starting with basic tokens and gradually adding complexity
3. **Pattern Recognition**: Identifying common patterns in language constructs
4. **Debugging Strategies**: Using print statements to understand token flow and AST structure

### Software Engineering Principles

1. **Separation of Concerns**: Clear distinction between lexing and parsing
2. **Interface Design**: Using Go interfaces for polymorphic AST nodes
3. **Code Organization**: Logical grouping of related functionality
4. **Documentation**: Extensive comments explaining the "why" behind each piece

### Language Design Insights

1. **Grammar Complexity**: Even "simple" languages have intricate rules
2. **Ambiguity Resolution**: How parsers handle potentially ambiguous syntax
3. **Error Recovery**: The challenge of meaningful error reporting
4. **Performance Considerations**: The trade-offs between simplicity and efficiency

This project gave me a deep appreciation for the complexity hidden in programming language tools we use every day. Every time I write JavaScript now, I have a better understanding of how the computer actually interprets my code!

---

_Built with ❤️ as a learning project to understand compiler theory and Go programming._
