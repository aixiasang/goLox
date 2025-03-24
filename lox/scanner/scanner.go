package scanner

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/token"
)

// Scanner 词法分析器，将源代码转换为标记列表
type Scanner struct {
	source  string         // 源代码
	tokens  []*token.Token // 标记列表
	start   int            // 当前词素的起始位置
	current int            // 当前字符的位置
	line    int            // 当前行号
	errors  error.Reporter // 错误报告器
	debug   bool           // 调试模式标志
}

// 关键字映射表
var keywords = map[string]token.TokenType{
	"and":    token.AND,
	"break":  token.BREAK,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

// NewScanner 创建一个新的词法分析器
func NewScanner(source string, errors error.Reporter) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []*token.Token{},
		start:   0,
		current: 0,
		line:    1,
		errors:  errors,
		debug:   false, // 默认关闭调试模式
	}
}

// SetDebug 设置调试模式
func (s *Scanner) SetDebug(debug bool) {
	s.debug = debug
}

// 调试输出辅助函数
func (s *Scanner) debugPrintf(format string, args ...interface{}) {
	if s.debug {
		fmt.Printf(format, args...)
	}
}

// ScanTokens 扫描所有标记
func (s *Scanner) ScanTokens() []*token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens
}

// scanToken 扫描单个标记
func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	// 单字符标记
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.debugPrintf("发现逗号标记，行: %d\n", s.line)
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '%':
		s.addToken(token.MODULO) // 取模运算符
	case '?':
		s.addToken(token.QUESTION) // 三元运算符问号
	case ':':
		s.addToken(token.COLON) // 三元运算符冒号

	// 一个或两个字符的标记
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}

	// 注释和长标记
	case '/':
		if s.match('/') {
			// 单行注释，一直读到行尾
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			// 块注释，一直读到 */
			s.blockComment()
		} else {
			s.addToken(token.SLASH)
		}

	// 忽略空白字符
	case ' ', '\r', '\t':
		// 忽略空白
	case '\n':
		s.line++

	// 字符串字面量
	case '"':
		s.string()

	default:
		if unicode.IsDigit(rune(c)) {
			s.number()
		} else if unicode.IsLetter(rune(c)) || c == '_' {
			s.identifier()
		} else {
			s.errors.ReportError(s.line, "未识别的字符。")
		}
	}
}

// blockComment 处理块注释
func (s *Scanner) blockComment() {
	// 嵌套层数
	nesting := 1

	for nesting > 0 && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		} else if s.peek() == '/' && s.peekNext() == '*' {
			s.advance() // 跳过 /
			s.advance() // 跳过 *
			nesting++
			continue
		} else if s.peek() == '*' && s.peekNext() == '/' {
			s.advance() // 跳过 *
			s.advance() // 跳过 /
			nesting--
			continue
		}
		s.advance()
	}

	if nesting > 0 {
		s.errors.ReportError(s.line, "未闭合的块注释。")
	}
}

// string 处理字符串字面量
func (s *Scanner) string() {
	// 读取直到找到闭合的引号
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errors.ReportError(s.line, "未闭合的字符串。")
		return
	}

	// 闭合的引号
	s.advance()

	// 提取字符串值（去除引号）
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(token.STRING, value)
}

// number 处理数字字面量
func (s *Scanner) number() {
	// 读取整数部分
	for unicode.IsDigit(rune(s.peek())) {
		s.advance()
	}

	// 处理小数部分
	if s.peek() == '.' && unicode.IsDigit(rune(s.peekNext())) {
		// 消费小数点
		s.advance()

		// 读取小数部分
		for unicode.IsDigit(rune(s.peek())) {
			s.advance()
		}
	}

	// 转换为浮点数
	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.errors.ReportError(s.line, "无效的数字。")
		return
	}
	s.addTokenWithLiteral(token.NUMBER, value)
}

// identifier 处理标识符和关键字
func (s *Scanner) identifier() {
	// 读取标识符剩余部分
	for unicode.IsLetter(rune(s.peek())) || unicode.IsDigit(rune(s.peek())) || s.peek() == '_' {
		s.advance()
	}

	// 获取标识符文本
	text := s.source[s.start:s.current]

	// 检查是否是关键字
	tokenType, isKeyword := keywords[text]
	if !isKeyword {
		tokenType = token.IDENTIFIER
	}

	s.addToken(tokenType)
}

// 辅助方法

// isAtEnd 检查是否已经到达源代码的末尾
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// advance 消费当前字符并前进
func (s *Scanner) advance() byte {
	result := s.source[s.current]
	s.current++
	return result
}

// match 检查当前字符是否匹配预期字符，如果匹配则消费
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

// peek 查看当前字符但不消费
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

// peekNext 查看下一个字符但不消费
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

// addToken 添加没有字面量的标记
func (s *Scanner) addToken(tokenType token.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

// addTokenWithLiteral 添加带有字面量的标记
func (s *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, text, literal, s.line))
}
