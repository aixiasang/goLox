package parser

import (
	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/token"
)

// Parser 实现递归下降解析器
type Parser struct {
	tokens  []*token.Token // 标记列表
	current int            // 当前标记索引
	errors  error.Reporter // 错误报告器
}

// NewParser 创建一个新的解析器
func NewParser(tokens []*token.Token, errors error.Reporter) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		errors:  errors,
	}
}

// Parse 解析表达式
func (p *Parser) Parse() ast.Expr {
	defer p.handlePanic()
	return p.expression()
}

// expression → equality
func (p *Parser) expression() ast.Expr {
	return p.equality()
}

// equality → comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr
}

// term → factor ( ( "-" | "+" ) factor )*
func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr
}

// factor → unary ( ( "/" | "*" ) unary )*
func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr
}

// unary → ( "!" | "-" ) unary | primary
func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnary(operator, right)
	}

	return p.primary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return ast.NewLiteral(false)
	}
	if p.match(token.TRUE) {
		return ast.NewLiteral(true)
	}
	if p.match(token.NIL) {
		return ast.NewLiteral(nil)
	}

	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteral(p.previous().Literal)
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "期待')'。")
		return ast.NewGrouping(expr)
	}

	panic(p.error(p.peek(), "期待表达式。"))
}

// 辅助方法

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(t token.TokenType, message string) *token.Token {
	if p.check(t) {
		return p.advance()
	}

	panic(p.error(p.peek(), message))
}

// 错误处理

func (p *Parser) error(token *token.Token, message string) error.ParseError {
	p.errors.Error(token, message)
	return error.NewParseError()
}

func (p *Parser) handlePanic() {
	if r := recover(); r != nil {
		if _, ok := r.(error.ParseError); !ok {
			panic(r)
		}
	}
}

// synchronize 在发生错误后重新同步解析器
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR,
			token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
