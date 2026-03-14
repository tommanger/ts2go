package types

import (
	"fmt"

	"github.com/tmanger/ts2go/internal/ast"
)

// TypeMapper handles TypeScript to Go type conversions
type TypeMapper struct {
	typeMap map[string]string
}

func NewTypeMapper() *TypeMapper {
	return &TypeMapper{
		typeMap: map[string]string{
			"string":  "string",
			"number":  "float64",
			"boolean": "bool",
			"any":     "interface{}",
			"void":    "",
		},
	}
}

// MapType converts a TypeScript type to a Go type
func (tm *TypeMapper) MapType(typeNode ast.TypeNode) string {
	if typeNode == nil {
		return ""
	}

	switch t := typeNode.(type) {
	case *ast.NamedType:
		if goType, ok := tm.typeMap[t.Name]; ok {
			return goType
		}
		// Custom type - keep the name (could be an interface or struct)
		return t.Name

	case *ast.ArrayType:
		elementType := tm.MapType(t.ElementType)
		return fmt.Sprintf("[]%s", elementType)

	case *ast.FunctionType:
		// For now, represent function types as interface{}
		// A more sophisticated implementation would create proper function signatures
		return "interface{}"

	default:
		return "interface{}"
	}
}

// InferType tries to infer the Go type from an expression
func (tm *TypeMapper) InferType(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return "float64"
	case *ast.StringLiteral:
		return "string"
	case *ast.BooleanLiteral:
		return "bool"
	case *ast.ArrayLiteral:
		if len(e.Elements) > 0 {
			elementType := tm.InferType(e.Elements[0])
			return fmt.Sprintf("[]%s", elementType)
		}
		return "[]interface{}"
	case *ast.BinaryExpression:
		leftType := tm.InferType(e.Left)
		rightType := tm.InferType(e.Right)
		// Simple heuristic: prefer the left type
		if leftType != "" {
			return leftType
		}
		return rightType
	default:
		return ""
	}
}
