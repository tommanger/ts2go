package parser

import (
	"fmt"

	"github.com/tmanger/ts2go/internal/ast"
	"github.com/tmanger/ts2go/internal/lexer"
)

type Parser struct {
	tokens  []lexer.Token
	pos     int
	current lexer.Token
}

func New(tokens []lexer.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    0,
	}
	if len(tokens) > 0 {
		p.current = tokens[0]
	}
	return p
}

func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.current.Type != lexer.TOKEN_EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program, nil
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens)-1 {
		p.pos++
		p.current = p.tokens[p.pos]
	}
}

func (p *Parser) peek() lexer.Token {
	if p.pos < len(p.tokens)-1 {
		return p.tokens[p.pos+1]
	}
	return lexer.Token{Type: lexer.TOKEN_EOF}
}

func (p *Parser) expect(tokenType lexer.TokenType) error {
	if p.current.Type != tokenType {
		return fmt.Errorf("expected %s, got %s at line %d column %d",
			tokenType, p.current.Type, p.current.Line, p.current.Column)
	}
	p.advance()
	return nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	// Skip semicolons
	if p.current.Type == lexer.TOKEN_SEMICOLON {
		p.advance()
		return nil, nil
	}

	var exported bool
	if p.current.Type == lexer.TOKEN_EXPORT {
		exported = true
		p.advance()
	}

	switch p.current.Type {
	case lexer.TOKEN_FUNCTION:
		return p.parseFunctionDeclaration(exported)
	case lexer.TOKEN_CONST, lexer.TOKEN_LET, lexer.TOKEN_VAR:
		return p.parseVariableDeclaration()
	case lexer.TOKEN_RETURN:
		return p.parseReturnStatement()
	case lexer.TOKEN_IF:
		return p.parseIfStatement()
	case lexer.TOKEN_FOR:
		return p.parseForStatement()
	case lexer.TOKEN_WHILE:
		return p.parseWhileStatement()
	case lexer.TOKEN_LBRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFunctionDeclaration(exported bool) (*ast.FunctionDeclaration, error) {
	p.advance() // skip 'function'

	if p.current.Type != lexer.TOKEN_IDENT {
		return nil, fmt.Errorf("expected function name at line %d", p.current.Line)
	}
	name := p.current.Literal
	p.advance()

	if err := p.expect(lexer.TOKEN_LPAREN); err != nil {
		return nil, err
	}

	params, err := p.parseParameters()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.TOKEN_RPAREN); err != nil {
		return nil, err
	}

	var returnType ast.TypeNode
	if p.current.Type == lexer.TOKEN_COLON {
		p.advance()
		returnType, err = p.parseType()
		if err != nil {
			return nil, err
		}
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionDeclaration{
		Name:       name,
		Params:     params,
		ReturnType: returnType,
		Body:       body,
		Exported:   exported,
	}, nil
}

func (p *Parser) parseParameters() ([]ast.Parameter, error) {
	params := []ast.Parameter{}

	if p.current.Type == lexer.TOKEN_RPAREN {
		return params, nil
	}

	for {
		if p.current.Type != lexer.TOKEN_IDENT {
			return nil, fmt.Errorf("expected parameter name at line %d", p.current.Line)
		}

		param := ast.Parameter{
			Name: p.current.Literal,
		}
		p.advance()

		if p.current.Type == lexer.TOKEN_COLON {
			p.advance()
			typeNode, err := p.parseType()
			if err != nil {
				return nil, err
			}
			param.Type = typeNode
		}

		params = append(params, param)

		if p.current.Type != lexer.TOKEN_COMMA {
			break
		}
		p.advance()
	}

	return params, nil
}

func (p *Parser) parseType() (ast.TypeNode, error) {
	if p.current.Type != lexer.TOKEN_IDENT {
		return nil, fmt.Errorf("expected type name at line %d", p.current.Line)
	}

	typeName := p.current.Literal
	p.advance()

	// Check for array type
	if p.current.Type == lexer.TOKEN_LBRACKET {
		p.advance()
		if err := p.expect(lexer.TOKEN_RBRACKET); err != nil {
			return nil, err
		}
		return &ast.ArrayType{
			ElementType: &ast.NamedType{Name: typeName},
		}, nil
	}

	return &ast.NamedType{Name: typeName}, nil
}

func (p *Parser) parseVariableDeclaration() (*ast.VariableDeclaration, error) {
	kind := p.current.Literal
	p.advance()

	if p.current.Type != lexer.TOKEN_IDENT {
		return nil, fmt.Errorf("expected variable name at line %d", p.current.Line)
	}

	name := p.current.Literal
	p.advance()

	var typeNode ast.TypeNode
	if p.current.Type == lexer.TOKEN_COLON {
		p.advance()
		var err error
		typeNode, err = p.parseType()
		if err != nil {
			return nil, err
		}
	}

	var value ast.Expression
	if p.current.Type == lexer.TOKEN_ASSIGN {
		p.advance()
		var err error
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	// Skip optional semicolon
	if p.current.Type == lexer.TOKEN_SEMICOLON {
		p.advance()
	}

	return &ast.VariableDeclaration{
		Kind:  kind,
		Name:  name,
		Type:  typeNode,
		Value: value,
	}, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	if err := p.expect(lexer.TOKEN_LBRACE); err != nil {
		return nil, err
	}

	block := &ast.BlockStatement{
		Statements: []ast.Statement{},
	}

	for p.current.Type != lexer.TOKEN_RBRACE && p.current.Type != lexer.TOKEN_EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	if err := p.expect(lexer.TOKEN_RBRACE); err != nil {
		return nil, err
	}

	return block, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	p.advance() // skip 'return'

	var value ast.Expression
	if p.current.Type != lexer.TOKEN_SEMICOLON && p.current.Type != lexer.TOKEN_RBRACE {
		var err error
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	if p.current.Type == lexer.TOKEN_SEMICOLON {
		p.advance()
	}

	return &ast.ReturnStatement{Value: value}, nil
}

func (p *Parser) parseIfStatement() (*ast.IfStatement, error) {
	p.advance() // skip 'if'

	if err := p.expect(lexer.TOKEN_LPAREN); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.TOKEN_RPAREN); err != nil {
		return nil, err
	}

	thenBlock, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	var elseBlock *ast.BlockStatement
	if p.current.Type == lexer.TOKEN_ELSE {
		p.advance()
		elseBlock, err = p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
	}

	return &ast.IfStatement{
		Condition: condition,
		ThenBlock: thenBlock,
		ElseBlock: elseBlock,
	}, nil
}

func (p *Parser) parseForStatement() (*ast.ForStatement, error) {
	p.advance() // skip 'for'

	if err := p.expect(lexer.TOKEN_LPAREN); err != nil {
		return nil, err
	}

	var init ast.Statement
	if p.current.Type != lexer.TOKEN_SEMICOLON {
		var err error
		init, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	} else {
		p.advance()
	}

	var condition ast.Expression
	if p.current.Type != lexer.TOKEN_SEMICOLON {
		var err error
		condition, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	if err := p.expect(lexer.TOKEN_SEMICOLON); err != nil {
		return nil, err
	}

	var update ast.Expression
	if p.current.Type != lexer.TOKEN_RPAREN {
		var err error
		update, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	if err := p.expect(lexer.TOKEN_RPAREN); err != nil {
		return nil, err
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.ForStatement{
		Init:      init,
		Condition: condition,
		Update:    update,
		Body:      body,
	}, nil
}

func (p *Parser) parseWhileStatement() (*ast.WhileStatement, error) {
	p.advance() // skip 'while'

	if err := p.expect(lexer.TOKEN_LPAREN); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.TOKEN_RPAREN); err != nil {
		return nil, err
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.current.Type == lexer.TOKEN_SEMICOLON {
		p.advance()
	}

	return &ast.ExpressionStatement{Expression: expr}, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseAssignment()
}

func (p *Parser) parseAssignment() (ast.Expression, error) {
	expr, err := p.parseLogicalOr()
	if err != nil {
		return nil, err
	}

	if p.current.Type == lexer.TOKEN_ASSIGN {
		p.advance()
		right, err := p.parseAssignment()
		if err != nil {
			return nil, err
		}
		return &ast.AssignmentExpression{
			Left:  expr,
			Right: right,
		}, nil
	}

	return expr, nil
}

func (p *Parser) parseLogicalOr() (ast.Expression, error) {
	left, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	for p.current.Type == lexer.TOKEN_OR {
		op := p.current.Literal
		p.advance()
		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseLogicalAnd() (ast.Expression, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	for p.current.Type == lexer.TOKEN_AND {
		op := p.current.Literal
		p.advance()
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseEquality() (ast.Expression, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.current.Type == lexer.TOKEN_EQUAL || p.current.Type == lexer.TOKEN_NOT_EQUAL {
		op := p.current.Literal
		p.advance()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseComparison() (ast.Expression, error) {
	left, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}

	for p.current.Type == lexer.TOKEN_LESS || p.current.Type == lexer.TOKEN_LESS_EQUAL ||
		p.current.Type == lexer.TOKEN_GREATER || p.current.Type == lexer.TOKEN_GREATER_EQUAL {
		op := p.current.Literal
		p.advance()
		right, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseAdditive() (ast.Expression, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}

	for p.current.Type == lexer.TOKEN_PLUS || p.current.Type == lexer.TOKEN_MINUS {
		op := p.current.Literal
		p.advance()
		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseMultiplicative() (ast.Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.current.Type == lexer.TOKEN_STAR || p.current.Type == lexer.TOKEN_SLASH || p.current.Type == lexer.TOKEN_PERCENT {
		op := p.current.Literal
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseUnary() (ast.Expression, error) {
	if p.current.Type == lexer.TOKEN_NOT || p.current.Type == lexer.TOKEN_MINUS {
		op := p.current.Literal
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpression{
			Operator: op,
			Operand:  operand,
		}, nil
	}

	return p.parsePostfix()
}

func (p *Parser) parsePostfix() (ast.Expression, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		switch p.current.Type {
		case lexer.TOKEN_LPAREN:
			p.advance()
			args, err := p.parseArguments()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.TOKEN_RPAREN); err != nil {
				return nil, err
			}
			expr = &ast.CallExpression{
				Function:  expr,
				Arguments: args,
			}
		case lexer.TOKEN_DOT:
			p.advance()
			if p.current.Type != lexer.TOKEN_IDENT {
				return nil, fmt.Errorf("expected property name at line %d", p.current.Line)
			}
			property := p.current.Literal
			p.advance()
			expr = &ast.MemberExpression{
				Object:   expr,
				Property: property,
			}
		case lexer.TOKEN_LBRACKET:
			p.advance()
			index, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.TOKEN_RBRACKET); err != nil {
				return nil, err
			}
			// Represent array access as a special call expression
			expr = &ast.CallExpression{
				Function:  expr,
				Arguments: []ast.Expression{index},
			}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) parseArguments() ([]ast.Expression, error) {
	args := []ast.Expression{}

	if p.current.Type == lexer.TOKEN_RPAREN {
		return args, nil
	}

	for {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		if p.current.Type != lexer.TOKEN_COMMA {
			break
		}
		p.advance()
	}

	return args, nil
}

func (p *Parser) parsePrimary() (ast.Expression, error) {
	switch p.current.Type {
	case lexer.TOKEN_NUMBER:
		value := p.current.Literal
		p.advance()
		return &ast.NumberLiteral{Value: value}, nil

	case lexer.TOKEN_STRING:
		value := p.current.Literal
		p.advance()
		return &ast.StringLiteral{Value: value}, nil

	case lexer.TOKEN_TRUE:
		p.advance()
		return &ast.BooleanLiteral{Value: true}, nil

	case lexer.TOKEN_FALSE:
		p.advance()
		return &ast.BooleanLiteral{Value: false}, nil

	case lexer.TOKEN_IDENT:
		name := p.current.Literal
		p.advance()
		return &ast.Identifier{Name: name}, nil

	case lexer.TOKEN_LPAREN:
		p.advance()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(lexer.TOKEN_RPAREN); err != nil {
			return nil, err
		}
		return expr, nil

	case lexer.TOKEN_LBRACKET:
		return p.parseArrayLiteral()

	default:
		return nil, fmt.Errorf("unexpected token %s at line %d", p.current.Type, p.current.Line)
	}
}

func (p *Parser) parseArrayLiteral() (*ast.ArrayLiteral, error) {
	p.advance() // skip '['

	elements := []ast.Expression{}

	if p.current.Type == lexer.TOKEN_RBRACKET {
		p.advance()
		return &ast.ArrayLiteral{Elements: elements}, nil
	}

	for {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, expr)

		if p.current.Type != lexer.TOKEN_COMMA {
			break
		}
		p.advance()
	}

	if err := p.expect(lexer.TOKEN_RBRACKET); err != nil {
		return nil, err
	}

	return &ast.ArrayLiteral{Elements: elements}, nil
}
