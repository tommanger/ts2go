package codegen

import (
	"fmt"
	"strings"

	"github.com/tmanger/ts2go/internal/ast"
	"github.com/tmanger/ts2go/internal/types"
)

type Generator struct {
	typeMapper *types.TypeMapper
	indent     int
}

func New() *Generator {
	return &Generator{
		typeMapper: types.NewTypeMapper(),
		indent:     0,
	}
}

func (g *Generator) Generate(program *ast.Program) string {
	var sb strings.Builder

	sb.WriteString("package main\n\n")

	// Check if we need fmt package
	needsFmt := g.needsFmtPackage(program)
	if needsFmt {
		sb.WriteString("import \"fmt\"\n\n")
	}

	// Separate function declarations from other statements
	var functions []*ast.FunctionDeclaration
	var topLevelStmts []ast.Statement

	for _, stmt := range program.Statements {
		if fn, ok := stmt.(*ast.FunctionDeclaration); ok {
			functions = append(functions, fn)
		} else {
			topLevelStmts = append(topLevelStmts, stmt)
		}
	}

	// Generate function declarations
	for _, fn := range functions {
		sb.WriteString(g.generateFunctionDeclaration(fn))
		sb.WriteString("\n\n")
	}

	// Generate main function if there are top-level statements
	if len(topLevelStmts) > 0 {
		sb.WriteString("func main() {\n")
		g.indent++
		for _, stmt := range topLevelStmts {
			sb.WriteString(g.generateStatement(stmt))
			sb.WriteString("\n")
		}
		g.indent--
		sb.WriteString("}\n")
	}

	return sb.String()
}

func (g *Generator) needsFmtPackage(program *ast.Program) bool {
	// Simple heuristic: check for console.log calls
	// In a real implementation, we'd do a proper AST traversal
	return true
}

func (g *Generator) generateStatement(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case *ast.FunctionDeclaration:
		return g.generateFunctionDeclaration(s)
	case *ast.VariableDeclaration:
		return g.generateVariableDeclaration(s)
	case *ast.BlockStatement:
		return g.generateBlockStatement(s)
	case *ast.ReturnStatement:
		return g.generateReturnStatement(s)
	case *ast.IfStatement:
		return g.generateIfStatement(s)
	case *ast.ForStatement:
		return g.generateForStatement(s)
	case *ast.WhileStatement:
		return g.generateWhileStatement(s)
	case *ast.ExpressionStatement:
		return g.indentation() + g.generateExpression(s.Expression)
	default:
		return ""
	}
}

func (g *Generator) generateFunctionDeclaration(fn *ast.FunctionDeclaration) string {
	var sb strings.Builder

	sb.WriteString(g.indentation())
	sb.WriteString("func ")

	// Capitalize first letter for exported functions
	name := fn.Name
	if fn.Exported {
		name = strings.ToUpper(name[:1]) + name[1:]
	}
	sb.WriteString(name)

	sb.WriteString("(")
	sb.WriteString(g.generateParamList(fn.Params))
	sb.WriteString(")")

	// Return type
	if fn.ReturnType != nil {
		returnType := g.typeMapper.MapType(fn.ReturnType)
		if returnType != "" {
			sb.WriteString(" ")
			sb.WriteString(returnType)
		}
	}

	sb.WriteString(" ")
	sb.WriteString(g.generateBlockStatement(fn.Body))

	return sb.String()
}

func (g *Generator) generateVariableDeclaration(v *ast.VariableDeclaration) string {
	var sb strings.Builder

	sb.WriteString(g.indentation())

	// If there's an explicit type annotation, use var declaration
	if v.Type != nil {
		sb.WriteString("var ")
		sb.WriteString(v.Name)
		sb.WriteString(" ")
		varType := g.typeMapper.MapType(v.Type)
		if varType == "" {
			varType = "interface{}"
		}
		sb.WriteString(varType)
		if v.Value != nil {
			sb.WriteString(" = ")
			sb.WriteString(g.generateExpression(v.Value))
		}
	} else if v.Value != nil {
		// Short variable declaration without type annotation
		sb.WriteString(v.Name)
		sb.WriteString(" := ")
		sb.WriteString(g.generateExpression(v.Value))
	} else {
		// Variable declaration without type or value (shouldn't happen in valid TS)
		sb.WriteString("var ")
		sb.WriteString(v.Name)
		sb.WriteString(" interface{}")
	}

	return sb.String()
}

func (g *Generator) generateBlockStatement(block *ast.BlockStatement) string {
	var sb strings.Builder

	sb.WriteString("{\n")
	g.indent++

	for _, stmt := range block.Statements {
		sb.WriteString(g.generateStatement(stmt))
		sb.WriteString("\n")
	}

	g.indent--
	sb.WriteString(g.indentation())
	sb.WriteString("}")

	return sb.String()
}

func (g *Generator) generateReturnStatement(r *ast.ReturnStatement) string {
	var sb strings.Builder

	sb.WriteString(g.indentation())
	sb.WriteString("return")

	if r.Value != nil {
		sb.WriteString(" ")
		sb.WriteString(g.generateExpression(r.Value))
	}

	return sb.String()
}

func (g *Generator) generateIfStatement(i *ast.IfStatement) string {
	var sb strings.Builder

	sb.WriteString(g.indentation())
	sb.WriteString("if ")
	sb.WriteString(g.generateExpression(i.Condition))
	sb.WriteString(" ")
	sb.WriteString(g.generateBlockStatement(i.ThenBlock))

	if i.ElseBlock != nil {
		sb.WriteString(" else ")
		sb.WriteString(g.generateBlockStatement(i.ElseBlock))
	}

	return sb.String()
}

func (g *Generator) generateForStatement(f *ast.ForStatement) string {
	var sb strings.Builder

	sb.WriteString(g.indentation())
	sb.WriteString("for ")

	// Init
	if f.Init != nil {
		initStr := g.generateStatement(f.Init)
		// Remove indentation and newline from init
		initStr = strings.TrimSpace(initStr)
		sb.WriteString(initStr)
	}
	sb.WriteString("; ")

	// Condition
	if f.Condition != nil {
		sb.WriteString(g.generateExpression(f.Condition))
	}
	sb.WriteString("; ")

	// Update
	if f.Update != nil {
		sb.WriteString(g.generateExpression(f.Update))
	}

	sb.WriteString(" ")
	sb.WriteString(g.generateBlockStatement(f.Body))

	return sb.String()
}

func (g *Generator) generateWhileStatement(w *ast.WhileStatement) string {
	var sb strings.Builder

	sb.WriteString(g.indentation())
	sb.WriteString("for ")
	sb.WriteString(g.generateExpression(w.Condition))
	sb.WriteString(" ")
	sb.WriteString(g.generateBlockStatement(w.Body))

	return sb.String()
}

func (g *Generator) generateExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return e.Name
	case *ast.NumberLiteral:
		return e.Value
	case *ast.StringLiteral:
		return fmt.Sprintf("\"%s\"", e.Value)
	case *ast.BooleanLiteral:
		if e.Value {
			return "true"
		}
		return "false"
	case *ast.ArrayLiteral:
		return g.generateArrayLiteral(e)
	case *ast.BinaryExpression:
		return g.generateBinaryExpression(e)
	case *ast.UnaryExpression:
		return g.generateUnaryExpression(e)
	case *ast.CallExpression:
		return g.generateCallExpression(e)
	case *ast.MemberExpression:
		return g.generateMemberExpression(e)
	case *ast.AssignmentExpression:
		return g.generateAssignmentExpression(e)
	case *ast.TemplateLiteral:
		return g.generateTemplateLiteral(e)
	case *ast.ArrowFunction:
		return g.generateArrowFunction(e)
	default:
		return ""
	}
}

func (g *Generator) generateArrayLiteral(a *ast.ArrayLiteral) string {
	var sb strings.Builder

	// Infer element type from first element if possible
	elementType := "interface{}"
	if len(a.Elements) > 0 {
		inferredType := g.typeMapper.InferType(a.Elements[0])
		if inferredType != "" {
			elementType = inferredType
		}
	}

	sb.WriteString("[]")
	sb.WriteString(elementType)
	sb.WriteString("{")

	for i, elem := range a.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(g.generateExpression(elem))
	}

	sb.WriteString("}")
	return sb.String()
}

func (g *Generator) generateBinaryExpression(b *ast.BinaryExpression) string {
	left := g.generateExpression(b.Left)
	right := g.generateExpression(b.Right)

	// Map TypeScript operators to Go operators
	op := b.Operator
	if op == "===" {
		op = "=="
	} else if op == "!==" {
		op = "!="
	}

	return fmt.Sprintf("(%s %s %s)", left, op, right)
}

func (g *Generator) generateUnaryExpression(u *ast.UnaryExpression) string {
	operand := g.generateExpression(u.Operand)
	return fmt.Sprintf("(%s%s)", u.Operator, operand)
}

func (g *Generator) generateCallExpression(c *ast.CallExpression) string {
	// Special handling for console.log -> fmt.Println
	if member, ok := c.Function.(*ast.MemberExpression); ok {
		if obj, ok := member.Object.(*ast.Identifier); ok && obj.Name == "console" {
			if member.Property == "log" {
				return g.generateFmtPrintln(c.Arguments)
			}
		}
	}

	// Regular function call
	var sb strings.Builder
	sb.WriteString(g.generateExpression(c.Function))
	sb.WriteString("(")

	for i, arg := range c.Arguments {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(g.generateExpression(arg))
	}

	sb.WriteString(")")
	return sb.String()
}

func (g *Generator) generateFmtPrintln(args []ast.Expression) string {
	var sb strings.Builder
	sb.WriteString("fmt.Println(")

	for i, arg := range args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(g.generateExpression(arg))
	}

	sb.WriteString(")")
	return sb.String()
}

func (g *Generator) generateMemberExpression(m *ast.MemberExpression) string {
	return fmt.Sprintf("%s.%s", g.generateExpression(m.Object), m.Property)
}

func (g *Generator) generateParamList(params []ast.Parameter) string {
	var sb strings.Builder
	for i, param := range params {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(param.Name)
		sb.WriteString(" ")
		paramType := g.typeMapper.MapType(param.Type)
		if paramType == "" {
			paramType = "interface{}"
		}
		sb.WriteString(paramType)
	}
	return sb.String()
}

func (g *Generator) generateArrowFunction(a *ast.ArrowFunction) string {
	var sb strings.Builder

	sb.WriteString("func(")
	sb.WriteString(g.generateParamList(a.Params))
	sb.WriteString(")")

	// Return type
	if a.ReturnType != nil {
		returnType := g.typeMapper.MapType(a.ReturnType)
		if returnType != "" {
			sb.WriteString(" ")
			sb.WriteString(returnType)
		}
	}

	sb.WriteString(" ")

	switch body := a.Body.(type) {
	case *ast.BlockStatement:
		sb.WriteString(g.generateBlockStatement(body))
	default:
		// Concise body: wrap expression in { return expr }
		sb.WriteString("{\n")
		g.indent++
		sb.WriteString(g.indentation())
		sb.WriteString("return ")
		sb.WriteString(g.generateExpression(body.(ast.Expression)))
		sb.WriteString("\n")
		g.indent--
		sb.WriteString(g.indentation())
		sb.WriteString("}")
	}

	return sb.String()
}

func (g *Generator) generateTemplateLiteral(t *ast.TemplateLiteral) string {
	if len(t.Expressions) == 0 {
		// No interpolations, just a plain string
		return fmt.Sprintf("\"%s\"", t.Parts[0])
	}

	// Build fmt.Sprintf format string
	var formatParts []string
	for _, part := range t.Parts {
		formatParts = append(formatParts, part)
	}

	// Join parts with %v placeholders
	var format strings.Builder
	for i, part := range formatParts {
		format.WriteString(part)
		if i < len(formatParts)-1 {
			format.WriteString("%v")
		}
	}

	// Generate expression arguments
	var args []string
	for _, expr := range t.Expressions {
		args = append(args, g.generateExpression(expr))
	}

	return fmt.Sprintf("fmt.Sprintf(\"%s\", %s)", format.String(), strings.Join(args, ", "))
}

func (g *Generator) generateAssignmentExpression(a *ast.AssignmentExpression) string {
	left := g.generateExpression(a.Left)
	right := g.generateExpression(a.Right)
	return fmt.Sprintf("%s = %s", left, right)
}

func (g *Generator) indentation() string {
	return strings.Repeat("\t", g.indent)
}
