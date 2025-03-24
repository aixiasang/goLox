package token

import "fmt"

// TokenType 表示词法分析器识别的不同类型的标记
type TokenType int

// Token类型常量定义
const (
	// 单字符标记
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// 一个或两个字符的标记
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// 字面量
	IDENTIFIER
	STRING
	NUMBER

	// 关键字
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

// Token名称映射表，用于调试和错误信息
var TokenNames = map[TokenType]string{
	LEFT_PAREN:    "LEFT_PAREN",
	RIGHT_PAREN:   "RIGHT_PAREN",
	LEFT_BRACE:    "LEFT_BRACE",
	RIGHT_BRACE:   "RIGHT_BRACE",
	COMMA:         "COMMA",
	DOT:           "DOT",
	MINUS:         "MINUS",
	PLUS:          "PLUS",
	SEMICOLON:     "SEMICOLON",
	SLASH:         "SLASH",
	STAR:          "STAR",
	BANG:          "BANG",
	BANG_EQUAL:    "BANG_EQUAL",
	EQUAL:         "EQUAL",
	EQUAL_EQUAL:   "EQUAL_EQUAL",
	GREATER:       "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS:          "LESS",
	LESS_EQUAL:    "LESS_EQUAL",
	IDENTIFIER:    "IDENTIFIER",
	STRING:        "STRING",
	NUMBER:        "NUMBER",
	AND:           "AND",
	CLASS:         "CLASS",
	ELSE:          "ELSE",
	FALSE:         "FALSE",
	FUN:           "FUN",
	FOR:           "FOR",
	IF:            "IF",
	NIL:           "NIL",
	OR:            "OR",
	PRINT:         "PRINT",
	RETURN:        "RETURN",
	SUPER:         "SUPER",
	THIS:          "THIS",
	TRUE:          "TRUE",
	VAR:           "VAR",
	WHILE:         "WHILE",
	EOF:           "EOF",
}

// Token 表示源代码中的词法标记
type Token struct {
	Type    TokenType
	Lexeme  string      // 原始词素文本
	Literal interface{} // 字面量值（对字符串、数字等）
	Line    int         // 行号，用于错误报告
}

// NewToken 创建一个新的Token实例
func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

// GetTokenName 返回TokenType对应的名称
func GetTokenName(tokenType TokenType) string {
	if name, exists := TokenNames[tokenType]; exists {
		return name
	}
	return "UNKNOWN"
}

// String 返回Token的字符串表示
func (t *Token) String() string {
	literalStr := ""
	if t.Literal != nil {
		literalStr = fmt.Sprintf(" %v", t.Literal)
	}
	return fmt.Sprintf("%s '%s'%s", GetTokenName(t.Type), t.Lexeme, literalStr)
}
