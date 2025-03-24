package lox

import (
	"fmt"
	"strconv"
)

// Scanner 是Lox语言的词法分析器
type Scanner struct {
	source  string   // 源代码
	tokens  []*Token // 扫描产生的token列表
	start   int      // 当前词素的起始位置
	current int      // 当前正在检查的字符位置
	line    int      // 当前行号
	lox     *Lox     // 用于报告错误

	// 保留字映射表
	keywords map[string]TokenType
}

// NewScanner 创建一个新的词法分析器
func NewScanner(source string) *Scanner {
	scanner := &Scanner{
		source:   source,
		tokens:   make([]*Token, 0),
		start:    0,
		current:  0,
		line:     1,
		lox:      NewLox(),
		keywords: make(map[string]TokenType),
	}

	// 初始化保留字映射表
	scanner.initKeywords()

	return scanner
}

// 初始化保留字映射表
func (s *Scanner) initKeywords() {
	s.keywords["and"] = AND
	s.keywords["class"] = CLASS
	s.keywords["else"] = ELSE
	s.keywords["false"] = FALSE
	s.keywords["for"] = FOR
	s.keywords["fun"] = FUN
	s.keywords["if"] = IF
	s.keywords["nil"] = NIL
	s.keywords["or"] = OR
	s.keywords["print"] = PRINT
	s.keywords["return"] = RETURN
	s.keywords["super"] = SUPER
	s.keywords["this"] = THIS
	s.keywords["true"] = TRUE
	s.keywords["var"] = VAR
	s.keywords["while"] = WHILE
}

// ScanTokens 扫描源代码并返回所有标记
func (s *Scanner) ScanTokens() []*Token {
	for !s.isAtEnd() {
		// 我们处于每个词素的开始位置
		s.start = s.current
		s.scanToken()
	}

	// 在文件末尾添加EOF标记
	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

// scanToken 扫描并识别单个标记
func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	// 单字符标记
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)

	// 一个或两个字符的标记
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}

	// 注释和除法
	case '/':
		if s.match('/') {
			// 单行注释，一直读到行尾
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			// 多行注释，支持嵌套
			s.multilineComment()
		} else {
			s.addToken(SLASH)
		}

	// 忽略空白字符
	case ' ', '\r', '\t':
		// 忽略
	case '\n':
		s.line++

	// 字符串字面量
	case '"':
		s.string()

	default:
		// 数字字面量
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			// 标识符和关键字
			s.identifier()
		} else {
			// 未知字符，报告错误
			s.lox.Error(s.line, fmt.Sprintf("Unexpected character: %c", c))
		}
	}
}

// string 处理字符串字面量
func (s *Scanner) string() {
	// 一直读取，直到找到结束的引号
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	// 如果到达文件结尾，说明字符串没有正确结束
	if s.isAtEnd() {
		s.lox.Error(s.line, "Unterminated string.")
		return
	}

	// 消费结束引号
	s.advance()

	// 提取字符串值（去掉两边的引号）
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, value)
}

// number 处理数字字面量
func (s *Scanner) number() {
	// 读取整数部分
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// 处理小数部分
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// 消费小数点
		s.advance()

		// 读取小数部分
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	// 转换为数值
	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.lox.Error(s.line, "Invalid number.")
		return
	}

	s.addTokenWithLiteral(NUMBER, value)
}

// identifier 处理标识符和关键字
func (s *Scanner) identifier() {
	// 读取整个标识符
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	// 检查是否是关键字
	text := s.source[s.start:s.current]
	tokenType, isKeyword := s.keywords[text]
	if !isKeyword {
		tokenType = IDENTIFIER
	}

	s.addToken(tokenType)
}

// advance 前进一个字符并返回当前字符
func (s *Scanner) advance() byte {
	if s.isAtEnd() {
		return 0
	}
	s.current++
	return s.source[s.current-1]
}

// match 如果下一个字符匹配预期，则消费它并返回true
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

// peek 查看当前字符但不消费它
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

// peekNext 向前查看两个字符
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

// isDigit 检查字符是否是数字
func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isAlpha 检查字符是否是字母或下划线
func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isAlphaNumeric 检查字符是否是字母、数字或下划线
func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

// addToken 添加一个没有字面量值的标记
func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

// addTokenWithLiteral 添加一个带有字面量值的标记
func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

// isAtEnd 检查是否到达源代码末尾
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// multilineComment 处理多行注释，支持嵌套
func (s *Scanner) multilineComment() {
	// 记录嵌套深度
	nestLevel := 1

	for nestLevel > 0 && !s.isAtEnd() {
		if s.peek() == '/' && s.peekNext() == '*' {
			// 嵌套注释开始
			s.advance() // 消费 '/'
			s.advance() // 消费 '*'
			nestLevel++
		} else if s.peek() == '*' && s.peekNext() == '/' {
			// 注释结束
			s.advance() // 消费 '*'
			s.advance() // 消费 '/'
			nestLevel--
		} else {
			// 处理换行
			if s.peek() == '\n' {
				s.line++
			}
			s.advance()
		}
	}

	// 如果到达文件结尾但注释未关闭
	if s.isAtEnd() && nestLevel > 0 {
		s.lox.Error(s.line, "Unterminated comment.")
	}
}
