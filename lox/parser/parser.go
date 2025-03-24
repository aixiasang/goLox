package parser

import (
	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/token"
)

// Parser 解析器，将标记序列转换为表达式
type Parser struct {
	tokens        []*token.Token // 标记列表
	current       int            // 当前标记索引
	errorReporter error.Reporter // 错误报告器
}

// NewParser 创建一个新的解析器
func NewParser(tokens []*token.Token, errorReporter error.Reporter) *Parser {
	return &Parser{
		tokens:        tokens,
		current:       0,
		errorReporter: errorReporter,
	}
}

// Parse 解析标记流，生成语法树
func (p *Parser) Parse() ast.Expr {
	defer p.handlePanic()
	return p.expression()
}

// handlePanic 处理解析过程中的异常
func (p *Parser) handlePanic() {
	if r := recover(); r != nil {
		if _, ok := r.(error.ParseError); ok {
			// 已经报告了解析错误，继续执行
			return
		}
		// 重新抛出其他类型的异常
		panic(r)
	}
}

// expression 解析表达式
func (p *Parser) expression() ast.Expr {
	return p.comma()
}

// comma 解析逗号表达式
func (p *Parser) comma() ast.Expr {
	var exprs []ast.Expr
	exprs = append(exprs, p.conditional())

	for p.match(token.COMMA) {
		exprs = append(exprs, p.conditional())
	}

	// 如果只有一个表达式，直接返回
	if len(exprs) == 1 {
		return exprs[0]
	}

	// 从右到左构建二叉树
	expr := exprs[len(exprs)-1]
	for i := len(exprs) - 2; i >= 0; i-- {
		expr = ast.NewBinary(exprs[i], &token.Token{
			Type:    token.COMMA,
			Lexeme:  ",",
			Literal: nil,
			Line:    p.previous().Line,
		}, expr)
	}

	return expr
}

// conditional 解析条件表达式（三元运算符）
func (p *Parser) conditional() ast.Expr {
	expr := p.equality()

	if p.match(token.QUESTION) {
		thenBranch := p.expression()
		if !p.match(token.COLON) {
			p.error(p.peek(), "期望在条件表达式中的 '?' 后有 ':'")
		}
		elseBranch := p.conditional()
		expr = ast.NewTernary(expr, thenBranch, elseBranch)
	}

	return expr
}

// equality 解析相等性表达式
func (p *Parser) equality() ast.Expr {
	if p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) && p.check(token.NUMBER, token.STRING, token.TRUE, token.FALSE, token.NIL, token.IDENTIFIER, token.LEFT_PAREN) {
		operator := p.previous()
		p.error(operator, "二元运算符缺少左操作数")
		right := p.comparison()
		return ast.NewBinary(ast.NewLiteral(nil), operator, right)
	}

	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr
}

// comparison 解析比较表达式
func (p *Parser) comparison() ast.Expr {
	if p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) && p.check(token.NUMBER, token.STRING, token.TRUE, token.FALSE, token.NIL, token.IDENTIFIER, token.LEFT_PAREN) {
		operator := p.previous()
		p.error(operator, "二元运算符缺少左操作数")
		right := p.term()
		return ast.NewBinary(ast.NewLiteral(nil), operator, right)
	}

	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr
}

// term 解析项表达式
func (p *Parser) term() ast.Expr {
	if p.match(token.PLUS) && p.check(token.NUMBER, token.STRING, token.TRUE, token.FALSE, token.NIL, token.IDENTIFIER, token.LEFT_PAREN) {
		operator := p.previous()
		p.error(operator, "二元运算符缺少左操作数")
		right := p.factor()
		return ast.NewBinary(ast.NewLiteral(nil), operator, right)
	}

	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr
}

// factor 解析因子表达式
func (p *Parser) factor() ast.Expr {
	if p.match(token.SLASH, token.STAR) && p.check(token.NUMBER, token.STRING, token.TRUE, token.FALSE, token.NIL, token.IDENTIFIER, token.LEFT_PAREN) {
		operator := p.previous()
		p.error(operator, "二元运算符缺少左操作数")
		right := p.unary()
		return ast.NewBinary(ast.NewLiteral(nil), operator, right)
	}

	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr
}

// unary 解析一元表达式
func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnary(operator, right)
	}

	return p.primary()
}

// primary 解析基本表达式
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

	if p.match(token.IDENTIFIER) {
		return ast.NewVariable(p.previous())
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "期望在表达式后有 ')'")
		return ast.NewGrouping(expr)
	}

	// 遇到错误，尝试同步恢复
	p.error(p.peek(), "期望表达式")
	return nil
}

// 辅助方法

// match 检查当前标记是否匹配任何给定类型，如果匹配则消费
func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

// check 检查当前标记是否为指定类型
func (p *Parser) check(types ...token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	current := p.peek().Type
	for _, t := range types {
		if current == t {
			return true
		}
	}
	return false
}

// advance 消费当前标记并返回
func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd 检查是否已经到达标记流结尾
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

// peek 返回当前标记但不消费
func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

// previous 返回最后一个消费的标记
func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

// consume 消费当前标记，如果不匹配则报错
func (p *Parser) consume(tokenType token.TokenType, message string) *token.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	p.error(p.peek(), message)
	panic(error.ParseError{Token: p.peek(), Message: message})
}

// error 报告当前标记的错误
func (p *Parser) error(token *token.Token, message string) {
	p.errorReporter.Error(token, 0, message)
	panic(error.ParseError{Token: token, Message: message})
}

// synchronize 尝试在错误后同步解析器状态
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
