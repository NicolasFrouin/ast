package main

import (
	"fmt"
)

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
