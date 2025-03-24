package interpreter

import (
	"testing"

	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/token"
)

// 创建一个模拟的错误报告器
type MockErrorReporter struct {
	Errors []string
}

func (m *MockErrorReporter) Error(tok *token.Token, line int, message string) {
	m.Errors = append(m.Errors, message)
}

func (m *MockErrorReporter) ReportError(line int, message string) {
	m.Errors = append(m.Errors, message)
}

func (m *MockErrorReporter) ResetError() {
	m.Errors = nil
}

func (m *MockErrorReporter) HasError() bool {
	return len(m.Errors) > 0
}

func (m *MockErrorReporter) HasRuntimeError() bool {
	return len(m.Errors) > 0
}

func TestLiteralExpression(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	tests := []struct {
		expr     ast.Expr
		expected interface{}
	}{
		{ast.NewLiteral(123.0), 123.0},
		{ast.NewLiteral("hello"), "hello"},
		{ast.NewLiteral(true), true},
		{ast.NewLiteral(nil), nil},
	}

	for _, tt := range tests {
		got := interpreter.evaluate(tt.expr)
		if got != tt.expected {
			t.Errorf("evaluate() = %v, want %v", got, tt.expected)
		}
	}
}

func TestUnaryExpression(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	tests := []struct {
		expr     ast.Expr
		expected interface{}
	}{
		{ast.NewUnary(token.NewToken(token.MINUS, "-", nil, 1), ast.NewLiteral(123.0)), -123.0},
		{ast.NewUnary(token.NewToken(token.BANG, "!", nil, 1), ast.NewLiteral(true)), false},
		{ast.NewUnary(token.NewToken(token.BANG, "!", nil, 1), ast.NewLiteral(false)), true},
		{ast.NewUnary(token.NewToken(token.BANG, "!", nil, 1), ast.NewLiteral(nil)), true},
	}

	for _, tt := range tests {
		got := interpreter.evaluate(tt.expr)
		if got != tt.expected {
			t.Errorf("evaluate() = %v, want %v", got, tt.expected)
		}
	}
}

func TestBinaryExpression(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	tests := []struct {
		expr     ast.Expr
		expected interface{}
	}{
		// 加法
		{ast.NewBinary(ast.NewLiteral(1.0), token.NewToken(token.PLUS, "+", nil, 1), ast.NewLiteral(2.0)), 3.0},
		{ast.NewBinary(ast.NewLiteral("Hello, "), token.NewToken(token.PLUS, "+", nil, 1), ast.NewLiteral("World!")), "Hello, World!"},

		// 减法
		{ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.MINUS, "-", nil, 1), ast.NewLiteral(3.0)), 2.0},

		// 乘法
		{ast.NewBinary(ast.NewLiteral(4.0), token.NewToken(token.STAR, "*", nil, 1), ast.NewLiteral(2.0)), 8.0},

		// 除法
		{ast.NewBinary(ast.NewLiteral(8.0), token.NewToken(token.SLASH, "/", nil, 1), ast.NewLiteral(2.0)), 4.0},

		// 大于
		{ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.GREATER, ">", nil, 1), ast.NewLiteral(3.0)), true},
		{ast.NewBinary(ast.NewLiteral(3.0), token.NewToken(token.GREATER, ">", nil, 1), ast.NewLiteral(5.0)), false},

		// 大于等于
		{ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.GREATER_EQUAL, ">=", nil, 1), ast.NewLiteral(5.0)), true},
		{ast.NewBinary(ast.NewLiteral(3.0), token.NewToken(token.GREATER_EQUAL, ">=", nil, 1), ast.NewLiteral(5.0)), false},

		// 小于
		{ast.NewBinary(ast.NewLiteral(3.0), token.NewToken(token.LESS, "<", nil, 1), ast.NewLiteral(5.0)), true},
		{ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.LESS, "<", nil, 1), ast.NewLiteral(3.0)), false},

		// 小于等于
		{ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.LESS_EQUAL, "<=", nil, 1), ast.NewLiteral(5.0)), true},
		{ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.LESS_EQUAL, "<=", nil, 1), ast.NewLiteral(3.0)), false},

		// 相等
		{ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.EQUAL_EQUAL, "==", nil, 1), ast.NewLiteral(5.0)), true},
		{ast.NewBinary(ast.NewLiteral("abc"), token.NewToken(token.EQUAL_EQUAL, "==", nil, 1), ast.NewLiteral("abc")), true},
		{ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.EQUAL_EQUAL, "==", nil, 1), ast.NewLiteral(3.0)), false},

		// 不等
		{ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.BANG_EQUAL, "!=", nil, 1), ast.NewLiteral(3.0)), true},
		{ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.BANG_EQUAL, "!=", nil, 1), ast.NewLiteral(5.0)), false},

		// 逗号
		{ast.NewBinary(ast.NewLiteral(1.0), token.NewToken(token.COMMA, ",", nil, 1), ast.NewLiteral(2.0)), 2.0},
	}

	for _, tt := range tests {
		got := interpreter.evaluate(tt.expr)
		if got != tt.expected {
			t.Errorf("evaluate() = %v, want %v", got, tt.expected)
		}
	}
}

func TestTernaryExpression(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	tests := []struct {
		expr     ast.Expr
		expected interface{}
	}{
		// 条件为true
		{ast.NewTernary(ast.NewLiteral(true), ast.NewLiteral(1.0), ast.NewLiteral(2.0)), 1.0},
		// 条件为false
		{ast.NewTernary(ast.NewLiteral(false), ast.NewLiteral(1.0), ast.NewLiteral(2.0)), 2.0},
		// 嵌套三元表达式
		{ast.NewTernary(
			ast.NewLiteral(true),
			ast.NewTernary(ast.NewLiteral(false), ast.NewLiteral(1.0), ast.NewLiteral(2.0)),
			ast.NewLiteral(3.0),
		), 2.0},
	}

	for _, tt := range tests {
		got := interpreter.evaluate(tt.expr)
		if got != tt.expected {
			t.Errorf("evaluate() = %v, want %v", got, tt.expected)
		}
	}
}

func TestGroupingExpression(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	tests := []struct {
		expr     ast.Expr
		expected interface{}
	}{
		{ast.NewGrouping(ast.NewLiteral(123.0)), 123.0},
		{ast.NewGrouping(ast.NewBinary(ast.NewLiteral(1.0), token.NewToken(token.PLUS, "+", nil, 1), ast.NewLiteral(2.0))), 3.0},
	}

	for _, tt := range tests {
		got := interpreter.evaluate(tt.expr)
		if got != tt.expected {
			t.Errorf("evaluate() = %v, want %v", got, tt.expected)
		}
	}
}

func TestRuntimeErrors(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	// 一元运算符错误: -"hello"
	unaryExpr := ast.NewUnary(token.NewToken(token.MINUS, "-", nil, 1), ast.NewLiteral("hello"))
	interpreter.Interpret(unaryExpr)
	if len(errorReporter.Errors) == 0 {
		t.Errorf("期望一元运算符错误")
	}
	errorReporter.ResetError()

	// 二元运算符错误: "hello" - 5
	binaryExpr := ast.NewBinary(ast.NewLiteral("hello"), token.NewToken(token.MINUS, "-", nil, 1), ast.NewLiteral(5.0))
	interpreter.Interpret(binaryExpr)
	if len(errorReporter.Errors) == 0 {
		t.Errorf("期望二元运算符错误")
	}
	errorReporter.ResetError()

	// 类型错误: "hello" + 5
	typeErrorExpr := ast.NewBinary(ast.NewLiteral("hello"), token.NewToken(token.PLUS, "+", nil, 1), ast.NewLiteral(5.0))
	interpreter.Interpret(typeErrorExpr)
	if len(errorReporter.Errors) == 0 {
		t.Errorf("期望类型错误")
	}
	errorReporter.ResetError()

	// 除零错误: 5 / 0
	divZeroExpr := ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.SLASH, "/", nil, 1), ast.NewLiteral(0.0))
	interpreter.Interpret(divZeroExpr)
	if len(errorReporter.Errors) == 0 {
		t.Errorf("期望除零错误")
	}
	errorReporter.ResetError()

	// 变量未定义错误
	varExpr := ast.NewVariable(token.NewToken(token.IDENTIFIER, "x", nil, 1))
	interpreter.Interpret(varExpr)
	if len(errorReporter.Errors) == 0 {
		t.Errorf("期望变量未定义错误")
	}
}

func TestStringify(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	tests := []struct {
		value    interface{}
		expected string
	}{
		{nil, "nil"},
		{123.0, "123"},
		{123.45, "123.45"},
		{"hello", "hello"},
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		got := interpreter.stringify(tt.value)
		if got != tt.expected {
			t.Errorf("stringify() = %v, want %v", got, tt.expected)
		}
	}
}
