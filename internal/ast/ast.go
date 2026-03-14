package ast

// Node is the base interface for all AST nodes
type Node interface {
	node()
}

// Statement nodes
type Statement interface {
	Node
	statement()
}

// Expression nodes
type Expression interface {
	Node
	expression()
}

// Type nodes
type TypeNode interface {
	Node
	typeNode()
}

// Program is the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) node() {}

// Type nodes
type NamedType struct {
	Name string
}

func (n *NamedType) node()     {}
func (n *NamedType) typeNode() {}

type ArrayType struct {
	ElementType TypeNode
}

func (a *ArrayType) node()     {}
func (a *ArrayType) typeNode() {}

type FunctionType struct {
	Params     []Parameter
	ReturnType TypeNode
}

func (f *FunctionType) node()     {}
func (f *FunctionType) typeNode() {}

// Parameter represents a function parameter
type Parameter struct {
	Name string
	Type TypeNode
}

// FunctionDeclaration represents a function declaration
type FunctionDeclaration struct {
	Name       string
	Params     []Parameter
	ReturnType TypeNode
	Body       *BlockStatement
	Exported   bool
}

func (f *FunctionDeclaration) node()      {}
func (f *FunctionDeclaration) statement() {}

// VariableDeclaration represents const, let, or var
type VariableDeclaration struct {
	Kind  string // "const", "let", or "var"
	Name  string
	Type  TypeNode
	Value Expression
}

func (v *VariableDeclaration) node()      {}
func (v *VariableDeclaration) statement() {}

// BlockStatement represents a block of statements
type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) node()      {}
func (b *BlockStatement) statement() {}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Value Expression
}

func (r *ReturnStatement) node()      {}
func (r *ReturnStatement) statement() {}

// IfStatement represents an if-else statement
type IfStatement struct {
	Condition Expression
	ThenBlock *BlockStatement
	ElseBlock *BlockStatement // can be nil
}

func (i *IfStatement) node()      {}
func (i *IfStatement) statement() {}

// ForStatement represents a for loop
type ForStatement struct {
	Init      Statement // can be nil or VariableDeclaration
	Condition Expression
	Update    Expression
	Body      *BlockStatement
}

func (f *ForStatement) node()      {}
func (f *ForStatement) statement() {}

// WhileStatement represents a while loop
type WhileStatement struct {
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileStatement) node()      {}
func (w *WhileStatement) statement() {}

// ExpressionStatement wraps an expression as a statement
type ExpressionStatement struct {
	Expression Expression
}

func (e *ExpressionStatement) node()      {}
func (e *ExpressionStatement) statement() {}

// Expressions

// Identifier represents a variable reference
type Identifier struct {
	Name string
}

func (i *Identifier) node()       {}
func (i *Identifier) expression() {}

// NumberLiteral represents a numeric literal
type NumberLiteral struct {
	Value string
}

func (n *NumberLiteral) node()       {}
func (n *NumberLiteral) expression() {}

// StringLiteral represents a string literal
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) node()       {}
func (s *StringLiteral) expression() {}

// BooleanLiteral represents true or false
type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) node()       {}
func (b *BooleanLiteral) expression() {}

// ArrayLiteral represents an array literal
type ArrayLiteral struct {
	Elements []Expression
}

func (a *ArrayLiteral) node()       {}
func (a *ArrayLiteral) expression() {}

// BinaryExpression represents binary operations
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpression) node()       {}
func (b *BinaryExpression) expression() {}

// UnaryExpression represents unary operations
type UnaryExpression struct {
	Operator string
	Operand  Expression
}

func (u *UnaryExpression) node()       {}
func (u *UnaryExpression) expression() {}

// CallExpression represents a function call
type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) node()       {}
func (c *CallExpression) expression() {}

// MemberExpression represents property access
type MemberExpression struct {
	Object   Expression
	Property string
}

func (m *MemberExpression) node()       {}
func (m *MemberExpression) expression() {}

// AssignmentExpression represents an assignment
type AssignmentExpression struct {
	Left  Expression
	Right Expression
}

func (a *AssignmentExpression) node()       {}
func (a *AssignmentExpression) expression() {}

// TemplateLiteral represents a template literal with interpolations
type TemplateLiteral struct {
	Parts       []string     // text segments (len = len(Expressions) + 1)
	Expressions []Expression // interpolated expressions
}

func (t *TemplateLiteral) node()       {}
func (t *TemplateLiteral) expression() {}

// ArrowFunction represents an arrow function expression
type ArrowFunction struct {
	Params     []Parameter
	ReturnType TypeNode
	Body       Node // *BlockStatement or Expression (concise body)
}

func (a *ArrowFunction) node()       {}
func (a *ArrowFunction) expression() {}
