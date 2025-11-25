package lexer

import (
	"strings"
	"unicode"
)

type TokenType int

const (
	TOKEN_EOF TokenType = iota
	TOKEN_EOL
	TOKEN_IDENTIFIER
	TOKEN_NUMBER
	TOKEN_FLOAT
	TOKEN_STRING
	TOKEN_BOOL
	TOKEN_DATE
	TOKEN_TIME

	// Opérateurs
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_MULTIPLY
	TOKEN_DIVIDE
	TOKEN_ASSIGN
	TOKEN_EQUAL
	TOKEN_NOT_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_IN
	TOKEN_LIKE
	TOKEN_BETWEEN
	TOKEN_RARROW
	TOKEN_LARROW
	TOKEN_NOT

	// Délimiteurs
	TOKEN_LPAREN
	TOKEN_RPAREN
	// TOKEN_LBRACE
	// TOKEN_RBRACE
	TOKEN_LBRAKET
	TOKEN_RBRAKET
	TOKEN_SEMICOLON
	TOKEN_COLON
	TOKEN_COMMA
	TOKEN_DOT

	// Mots-clés
	TOKEN_IF
	TOKEN_ELSE
	TOKEN_WHILE
	TOKEN_FOR
	TOKEN_FOREACH
	TOKEN_FUNCTION
	TOKEN_RETURN
	TOKEN_LET
	TOKEN_TYPE
	TOKEN_RECORD
	TOKEN_ACTION
	TOKEN_START
	TOKEN_END
	TOKEN_DO
	TOKEN_STOP
	TOKEN_NUMBER_TYPE
	TOKEN_FLOAT_TYPE
	TOKEN_STRING_TYPE
	TOKEN_BOOL_TYPE
	TOKEN_DATE_TYPE
	TOKEN_TIME_TYPE
	TOKEN_ARRAY
	TOKEN_SELECT
	TOKEN_FROM
	TOKEN_WHERE
	TOKEN_RECURSIVE
	TOKEN_BROWSE
	TOKEN_CASE
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

type Lexer struct {
	input  string
	pos    int
	line   int
	column int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: TOKEN_EOF, Line: l.line, Column: l.column}
	}

	// Commentaires
	if l.input[l.pos] == '(' && l.peek() == '*' {
		l.skipComment()
	}

	ch := l.input[l.pos]

	// Identifiants et mots-clés
	if unicode.IsLetter(rune(ch)) || ch == '_' {
		return l.readIdentifier()
	}

	// Nombres
	if unicode.IsDigit(rune(ch)) {
		return l.readNumber()
	}

	// Chaînes de caractères
	if ch == '"' || ch == '\'' {
		return l.readString()
	}

	// Opérateurs et délimiteurs
	switch ch {
	case '\r':
		return l.createToken(TOKEN_EOL, "\r")
	case '+':
		return l.createToken(TOKEN_PLUS, "+")
	case '-':
		if l.peek() == '>' {
			l.consume()
			return l.createToken(TOKEN_RARROW, "->")
		}
		return l.createToken(TOKEN_MINUS, "-")
	case '*':
		return l.createToken(TOKEN_MULTIPLY, "*")
	case '/':
		return l.createToken(TOKEN_DIVIDE, "/")
	case '=':
		if l.peek() == '=' {
			l.consume()
			return l.createToken(TOKEN_EQUAL, "==")
		}
		return l.createToken(TOKEN_ASSIGN, "=")
	case '<':
		if l.peek() == '=' {
			l.consume()
			return l.createToken(TOKEN_LESS_EQUAL, "<=")
		}
		if l.peek() == '>' {
			l.consume()
			return l.createToken(TOKEN_NOT_EQUAL, "<>")
		}
		if l.peek() == '-' {
			l.consume()
			return l.createToken(TOKEN_LARROW, "<-")
		}
		return l.createToken(TOKEN_LESS, "<")
	case '>':
		if l.peek() == '=' {
			l.consume()
			return l.createToken(TOKEN_GREATER_EQUAL, ">=")
		}
		return l.createToken(TOKEN_GREATER, ">")
	case '!':
		if l.peek() == '=' {
			l.consume()
			return l.createToken(TOKEN_NOT_EQUAL, "!=")
		}
		return l.createToken(TOKEN_NOT, "!")
	case '[':
		if l.peek() == '=' {
			l.consume()
			return l.createToken(TOKEN_LBRAKET, "[")
		}
	case ']':
		if l.peek() == '=' {
			l.consume()
			return l.createToken(TOKEN_RBRAKET, "]")
		}
	case '(':
		return l.createToken(TOKEN_LPAREN, "(")
	case ')':
		return l.createToken(TOKEN_RPAREN, ")")
	case ';':
		return l.createToken(TOKEN_SEMICOLON, ";")
	case ',':
		return l.createToken(TOKEN_COMMA, ",")
	case '.':
		return l.createToken(TOKEN_DOT, ".")
	case ':':
		return l.createToken(TOKEN_DOT, ".")
	}

	// Token inconnu
	token := l.createToken(TOKEN_EOF, string(ch))
	l.consume()
	return token
}

func (l *Lexer) readIdentifier() Token {
	start := l.pos
	for l.pos < len(l.input) && (unicode.IsLetter(rune(l.input[l.pos])) ||
		unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '_') {
		l.consume()
	}

	value := l.input[start:l.pos]
	tokenType := l.lookupKeyword(value)

	return Token{
		Type:   tokenType,
		Value:  value,
		Line:   l.line,
		Column: l.column - len(value),
	}
}

func (l *Lexer) lookupKeyword(ident string) TokenType {
	switch strings.ToLower(ident) {
	case "if":
		return TOKEN_IF
	case "else":
		return TOKEN_ELSE
	case "while":
		return TOKEN_WHILE
	case "select":
		return TOKEN_SELECT
	case "case":
		return TOKEN_CASE
	case "for":
		return TOKEN_FOR
	case "function":
		return TOKEN_FUNCTION
	case "return":
		return TOKEN_RETURN
	case "let":
		return TOKEN_LET
	case "type":
		return TOKEN_TYPE
	case "record":
		return TOKEN_RECORD
	case "action":
		return TOKEN_ACTION
	case "start":
		return TOKEN_START
	case "end":
		return TOKEN_END
	case "do":
		return TOKEN_DO
	case "stop":
		return TOKEN_STOP
	case "number":
		return TOKEN_NUMBER_TYPE
	case "float":
		return TOKEN_FLOAT_TYPE
	case "string":
		return TOKEN_STRING_TYPE
	case "boolean":
		return TOKEN_BOOL_TYPE
	case "date":
		return TOKEN_DATE_TYPE
	case "time":
		return TOKEN_TIME_TYPE
	case "array":
		return TOKEN_ARRAY
	case "from":
		return TOKEN_FROM
	case "where":
		return TOKEN_WHERE
	case "recursive":
		return TOKEN_RECURSIVE
	case "browse":
		return TOKEN_BROWSE
	case "in":
		return TOKEN_IN
	case "like":
		return TOKEN_LIKE
	case "between":
		return TOKEN_BETWEEN
	case "not":
		return TOKEN_NOT
	default:
		return TOKEN_IDENTIFIER
	}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
		l.consume()
	}

	value := l.input[start:l.pos]
	return Token{
		Type:   TOKEN_NUMBER,
		Value:  value,
		Line:   l.line,
		Column: l.column - len(value),
	}
}

func (l *Lexer) readString() Token {
	l.consume() // Skip opening quote
	start := l.pos

	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		if l.input[l.pos] == '\n' {
			l.line++
			l.column = 1
		}
		l.consume()
	}

	value := l.input[start:l.pos]
	l.consume() // Skip closing quote

	return Token{
		Type:   TOKEN_STRING,
		Value:  value,
		Line:   l.line,
		Column: l.column - len(value) - 2,
	}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.consume()
		} else if ch == '\n' {
			l.line++
			l.column = 1
			l.consume()
		} else {
			break
		}
	}
}

func (l *Lexer) skipComment() {
	if l.pos >= len(l.input) {
		return
	}
	if l.input[l.pos] != '(' || l.peek() != '*' {
		return
	}
	l.consume() //Reads '*'
	l.consume() //Move the cursor to the next position
	for ch := l.input[l.pos]; l.pos < len(l.input) &&
		ch != '*' && l.peek() != ')'; ch = l.input[l.pos] {
		if ch == '\n' {
			l.line++
			l.column = 1
		}
		l.consume()
	}
	if l.peek() == ')' {
		l.consume() //Reads ')'
		l.consume() //Move the cursor to the next position
	}
}

func (l *Lexer) createToken(tokenType TokenType, value string) Token {
	token := Token{
		Type:   tokenType,
		Value:  value,
		Line:   l.line,
		Column: l.column,
	}
	l.consumeN(len(value))
	return token
}

func (l *Lexer) consume() {
	if l.pos < len(l.input) {
		l.pos++
		l.column++
	}
}

func (l *Lexer) consumeN(n int) {
	for i := 0; i < n && l.pos < len(l.input); i++ {
		l.consume()
	}
}

func (l *Lexer) peek() byte {
	if l.pos+1 < len(l.input) {
		return l.input[l.pos+1]
	}
	return 0
}

