package lexer

type TokenType int

const (
	// Special tokens
	TOKEN_EOF TokenType = iota
	TOKEN_ILLEGAL

	// Identifiers and literals
	TOKEN_IDENT
	TOKEN_NUMBER
	TOKEN_STRING
	TOKEN_TRUE
	TOKEN_FALSE

	// Keywords
	TOKEN_FUNCTION
	TOKEN_CONST
	TOKEN_LET
	TOKEN_VAR
	TOKEN_IF
	TOKEN_ELSE
	TOKEN_FOR
	TOKEN_WHILE
	TOKEN_RETURN
	TOKEN_INTERFACE
	TOKEN_TYPE
	TOKEN_EXPORT

	// Operators
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_STAR
	TOKEN_SLASH
	TOKEN_PERCENT
	TOKEN_ASSIGN
	TOKEN_EQUAL
	TOKEN_NOT_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_AND
	TOKEN_OR
	TOKEN_NOT

	// Delimiters
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_LBRACKET
	TOKEN_RBRACKET
	TOKEN_SEMICOLON
	TOKEN_COLON
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_ARROW
)

var keywords = map[string]TokenType{
	"function":  TOKEN_FUNCTION,
	"const":     TOKEN_CONST,
	"let":       TOKEN_LET,
	"var":       TOKEN_VAR,
	"if":        TOKEN_IF,
	"else":      TOKEN_ELSE,
	"for":       TOKEN_FOR,
	"while":     TOKEN_WHILE,
	"return":    TOKEN_RETURN,
	"true":      TOKEN_TRUE,
	"false":     TOKEN_FALSE,
	"interface": TOKEN_INTERFACE,
	"type":      TOKEN_TYPE,
	"export":    TOKEN_EXPORT,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func (t TokenType) String() string {
	names := map[TokenType]string{
		TOKEN_EOF:           "EOF",
		TOKEN_ILLEGAL:       "ILLEGAL",
		TOKEN_IDENT:         "IDENT",
		TOKEN_NUMBER:        "NUMBER",
		TOKEN_STRING:        "STRING",
		TOKEN_TRUE:          "TRUE",
		TOKEN_FALSE:         "FALSE",
		TOKEN_FUNCTION:      "FUNCTION",
		TOKEN_CONST:         "CONST",
		TOKEN_LET:           "LET",
		TOKEN_VAR:           "VAR",
		TOKEN_IF:            "IF",
		TOKEN_ELSE:          "ELSE",
		TOKEN_FOR:           "FOR",
		TOKEN_WHILE:         "WHILE",
		TOKEN_RETURN:        "RETURN",
		TOKEN_INTERFACE:     "INTERFACE",
		TOKEN_TYPE:          "TYPE",
		TOKEN_EXPORT:        "EXPORT",
		TOKEN_PLUS:          "+",
		TOKEN_MINUS:         "-",
		TOKEN_STAR:          "*",
		TOKEN_SLASH:         "/",
		TOKEN_PERCENT:       "%",
		TOKEN_ASSIGN:        "=",
		TOKEN_EQUAL:         "==",
		TOKEN_NOT_EQUAL:     "!=",
		TOKEN_LESS:          "<",
		TOKEN_LESS_EQUAL:    "<=",
		TOKEN_GREATER:       ">",
		TOKEN_GREATER_EQUAL: ">=",
		TOKEN_AND:           "&&",
		TOKEN_OR:            "||",
		TOKEN_NOT:           "!",
		TOKEN_LPAREN:        "(",
		TOKEN_RPAREN:        ")",
		TOKEN_LBRACE:        "{",
		TOKEN_RBRACE:        "}",
		TOKEN_LBRACKET:      "[",
		TOKEN_RBRACKET:      "]",
		TOKEN_SEMICOLON:     ";",
		TOKEN_COLON:         ":",
		TOKEN_COMMA:         ",",
		TOKEN_DOT:           ".",
		TOKEN_ARROW:         "=>",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return "UNKNOWN"
}
