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

func TestVariableDeclarationAndAssignment(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	// 声明变量
	varStmt := ast.NewVar(
		token.NewToken(token.IDENTIFIER, "x", nil, 1),
		ast.NewLiteral(42.0),
	)

	// 执行声明语句
	interpreter.execute(varStmt)

	// 读取变量
	varExpr := ast.NewVariable(token.NewToken(token.IDENTIFIER, "x", nil, 1))
	value := interpreter.evaluate(varExpr)

	if value != 42.0 {
		t.Errorf("Variable 'x' should be 42.0, got %v", value)
	}

	// 赋值给变量
	assignExpr := ast.NewAssign(
		token.NewToken(token.IDENTIFIER, "x", nil, 1),
		ast.NewLiteral(100.0),
	)

	assignValue := interpreter.evaluate(assignExpr)
	if assignValue != 100.0 {
		t.Errorf("Assignment should return 100.0, got %v", assignValue)
	}

	// 验证变量值被更新
	updatedValue := interpreter.evaluate(varExpr)
	if updatedValue != 100.0 {
		t.Errorf("After assignment, 'x' should be 100.0, got %v", updatedValue)
	}
}

func TestBlockStatement(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	// 全局变量
	globalVarStmt := ast.NewVar(
		token.NewToken(token.IDENTIFIER, "x", nil, 1),
		ast.NewLiteral(10.0),
	)

	// 块内变量
	blockVarStmt := ast.NewVar(
		token.NewToken(token.IDENTIFIER, "x", nil, 2),
		ast.NewLiteral(20.0),
	)

	// 块内的另一个变量
	blockVarY := ast.NewVar(
		token.NewToken(token.IDENTIFIER, "y", nil, 3),
		ast.NewLiteral(30.0),
	)

	// 创建块
	blockStmt := ast.NewBlock([]ast.Stmt{
		blockVarStmt,
		blockVarY,
	})

	// 执行全局变量声明
	interpreter.execute(globalVarStmt)

	// 执行块
	interpreter.execute(blockStmt)

	// 块执行后，x应该仍然是10.0（全局变量）
	xExpr := ast.NewVariable(token.NewToken(token.IDENTIFIER, "x", nil, 4))
	xValue := interpreter.evaluate(xExpr)

	if xValue != 10.0 {
		t.Errorf("After block execution, global 'x' should be 10.0, got %v", xValue)
	}

	// y变量不应该存在于全局作用域
	yExpr := ast.NewVariable(token.NewToken(token.IDENTIFIER, "y", nil, 5))

	// 使用匿名函数捕获panic
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when accessing 'y' outside of block")
			}
		}()

		interpreter.evaluate(yExpr)
	}()
}

func TestIfStatement(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	// 测试变量
	xToken := token.NewToken(token.IDENTIFIER, "x", nil, 1)
	interpreter.environment.Define(xToken.Lexeme, 0.0)

	// 创建if语句：if (true) x = 1; else x = 2;
	ifStmt := ast.NewIf(
		ast.NewLiteral(true),
		ast.NewExpression(ast.NewAssign(xToken, ast.NewLiteral(1.0))),
		ast.NewExpression(ast.NewAssign(xToken, ast.NewLiteral(2.0))),
	)

	// 执行if语句
	interpreter.execute(ifStmt)

	// 检查x的值，应该是1.0
	xExpr := ast.NewVariable(xToken)
	xValue := interpreter.evaluate(xExpr)

	if xValue != 1.0 {
		t.Errorf("After if statement with true condition, 'x' should be 1.0, got %v", xValue)
	}

	// 创建if语句：if (false) x = 3; else x = 4;
	ifStmt2 := ast.NewIf(
		ast.NewLiteral(false),
		ast.NewExpression(ast.NewAssign(xToken, ast.NewLiteral(3.0))),
		ast.NewExpression(ast.NewAssign(xToken, ast.NewLiteral(4.0))),
	)

	// 执行if语句
	interpreter.execute(ifStmt2)

	// 检查x的值，应该是4.0
	xValue = interpreter.evaluate(xExpr)

	if xValue != 4.0 {
		t.Errorf("After if statement with false condition, 'x' should be 4.0, got %v", xValue)
	}
}

func TestWhileStatement(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	// 测试变量
	iToken := token.NewToken(token.IDENTIFIER, "i", nil, 1)
	interpreter.environment.Define(iToken.Lexeme, 0.0)

	// 创建while语句：while (i < 5) i = i + 1;
	whileStmt := ast.NewWhile(
		ast.NewBinary(
			ast.NewVariable(iToken),
			token.NewToken(token.LESS, "<", nil, 1),
			ast.NewLiteral(5.0),
		),
		ast.NewExpression(ast.NewAssign(
			iToken,
			ast.NewBinary(
				ast.NewVariable(iToken),
				token.NewToken(token.PLUS, "+", nil, 1),
				ast.NewLiteral(1.0),
			),
		)),
	)

	// 执行while语句
	interpreter.execute(whileStmt)

	// 检查i的值，应该是5.0
	iExpr := ast.NewVariable(iToken)
	iValue := interpreter.evaluate(iExpr)

	if iValue != 5.0 {
		t.Errorf("After while loop, 'i' should be 5.0, got %v", iValue)
	}
}

func TestRuntimeErrors(t *testing.T) {
	errorReporter := &MockErrorReporter{}
	interpreter := NewInterpreter(errorReporter)

	// 一元运算符错误: -"hello"
	unaryExpr := ast.NewUnary(token.NewToken(token.MINUS, "-", nil, 1), ast.NewLiteral("hello"))

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for unary operator error")
			}
		}()
		interpreter.evaluate(unaryExpr)
	}()

	// 执行语句列表确保错误被捕获
	statements := []ast.Stmt{ast.NewExpression(unaryExpr)}
	interpreter.Interpret(statements)

	if len(errorReporter.Errors) == 0 {
		t.Errorf("期望一元运算符错误")
	}
	errorReporter.ResetError()

	// 二元运算符错误: "hello" - 5
	binaryExpr := ast.NewBinary(ast.NewLiteral("hello"), token.NewToken(token.MINUS, "-", nil, 1), ast.NewLiteral(5.0))
	statements = []ast.Stmt{ast.NewExpression(binaryExpr)}
	interpreter.Interpret(statements)

	if len(errorReporter.Errors) == 0 {
		t.Errorf("期望二元运算符错误")
	}
	errorReporter.ResetError()

	// 类型错误: "hello" + 5
	typeErrorExpr := ast.NewBinary(ast.NewLiteral("hello"), token.NewToken(token.PLUS, "+", nil, 1), ast.NewLiteral(5.0))
	statements = []ast.Stmt{ast.NewExpression(typeErrorExpr)}
	interpreter.Interpret(statements)

	if len(errorReporter.Errors) == 0 {
		t.Errorf("期望类型错误")
	}
	errorReporter.ResetError()

	// 除零错误: 5 / 0
	divZeroExpr := ast.NewBinary(ast.NewLiteral(5.0), token.NewToken(token.SLASH, "/", nil, 1), ast.NewLiteral(0.0))
	statements = []ast.Stmt{ast.NewExpression(divZeroExpr)}
	interpreter.Interpret(statements)

	if len(errorReporter.Errors) == 0 {
		t.Errorf("期望除零错误")
	}
	errorReporter.ResetError()

	// 变量未定义错误
	varExpr := ast.NewVariable(token.NewToken(token.IDENTIFIER, "x", nil, 1))
	statements = []ast.Stmt{ast.NewExpression(varExpr)}
	interpreter.Interpret(statements)

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
