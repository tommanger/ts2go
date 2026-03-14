package lexer

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	input   string
	pos     int  // current position
	readPos int  // next reading position
	ch      byte // current char
	line    int
	column  int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.column++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type == TOKEN_EOF {
			break
		}
		if tok.Type == TOKEN_ILLEGAL {
			return nil, fmt.Errorf("illegal token at line %d, column %d: %s", tok.Line, tok.Column, tok.Literal)
		}
	}

	return tokens, nil
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_EQUAL, Literal: "==", Line: tok.Line, Column: tok.Column}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: TOKEN_ARROW, Literal: "=>", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: TOKEN_ASSIGN, Literal: "=", Line: tok.Line, Column: tok.Column}
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_NOT_EQUAL, Literal: "!=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: TOKEN_NOT, Literal: "!", Line: tok.Line, Column: tok.Column}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_LESS_EQUAL, Literal: "<=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: TOKEN_LESS, Literal: "<", Line: tok.Line, Column: tok.Column}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_GREATER_EQUAL, Literal: ">=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: TOKEN_GREATER, Literal: ">", Line: tok.Line, Column: tok.Column}
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = Token{Type: TOKEN_AND, Literal: "&&", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: TOKEN_ILLEGAL, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = Token{Type: TOKEN_OR, Literal: "||", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: TOKEN_ILLEGAL, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case '+':
		tok = Token{Type: TOKEN_PLUS, Literal: "+", Line: tok.Line, Column: tok.Column}
	case '-':
		tok = Token{Type: TOKEN_MINUS, Literal: "-", Line: tok.Line, Column: tok.Column}
	case '*':
		tok = Token{Type: TOKEN_STAR, Literal: "*", Line: tok.Line, Column: tok.Column}
	case '/':
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken()
		}
		tok = Token{Type: TOKEN_SLASH, Literal: "/", Line: tok.Line, Column: tok.Column}
	case '%':
		tok = Token{Type: TOKEN_PERCENT, Literal: "%", Line: tok.Line, Column: tok.Column}
	case '(':
		tok = Token{Type: TOKEN_LPAREN, Literal: "(", Line: tok.Line, Column: tok.Column}
	case ')':
		tok = Token{Type: TOKEN_RPAREN, Literal: ")", Line: tok.Line, Column: tok.Column}
	case '{':
		tok = Token{Type: TOKEN_LBRACE, Literal: "{", Line: tok.Line, Column: tok.Column}
	case '}':
		tok = Token{Type: TOKEN_RBRACE, Literal: "}", Line: tok.Line, Column: tok.Column}
	case '[':
		tok = Token{Type: TOKEN_LBRACKET, Literal: "[", Line: tok.Line, Column: tok.Column}
	case ']':
		tok = Token{Type: TOKEN_RBRACKET, Literal: "]", Line: tok.Line, Column: tok.Column}
	case ';':
		tok = Token{Type: TOKEN_SEMICOLON, Literal: ";", Line: tok.Line, Column: tok.Column}
	case ':':
		tok = Token{Type: TOKEN_COLON, Literal: ":", Line: tok.Line, Column: tok.Column}
	case ',':
		tok = Token{Type: TOKEN_COMMA, Literal: ",", Line: tok.Line, Column: tok.Column}
	case '.':
		tok = Token{Type: TOKEN_DOT, Literal: ".", Line: tok.Line, Column: tok.Column}
	case '"', '\'', '`':
		tok.Type = TOKEN_STRING
		tok.Literal = l.readString(l.ch)
	case 0:
		tok = Token{Type: TOKEN_EOF, Literal: "", Line: tok.Line, Column: tok.Column}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupKeyword(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TOKEN_NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{Type: TOKEN_ILLEGAL, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComment() {
	l.readChar() // skip *
	for {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // skip *
			l.readChar() // skip /
			break
		}
		if l.ch == 0 {
			break
		}
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readString(quote byte) string {
	pos := l.pos + 1
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar() // skip escaped char
		}
	}
	return l.input[pos:l.pos]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_' || ch == '$'
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

func lookupKeyword(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENT
}
