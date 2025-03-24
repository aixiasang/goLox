package interpreter

import (
	"fmt"
	"math"

	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/environment"
	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/token"
)

// Interpreter 实现表达式求值和语句执行
type Interpreter struct {
	errorReporter error.Reporter
	environment   *environment.Environment
}

// NewInterpreter 创建一个新的解释器
func NewInterpreter(errorReporter error.Reporter) *Interpreter {
	return &Interpreter{
		errorReporter: errorReporter,
		environment:   environment.NewEnvironment(),
	}
}

// Interpret 解释执行语句列表
func (i *Interpreter) Interpret(statements []ast.Stmt) {
	defer i.handlePanic()

	for _, stmt := range statements {
		i.execute(stmt)
	}
}

// execute 执行一条语句
func (i *Interpreter) execute(stmt ast.Stmt) {
	stmt.Accept(i)
}

// executeBlock 在给定环境中执行语句块
func (i *Interpreter) executeBlock(statements []ast.Stmt, env *environment.Environment) {
	previous := i.environment

	// 恢复原来的环境，即使panic也要确保恢复
	defer func() {
		i.environment = previous
	}()

	i.environment = env

	for _, statement := range statements {
		i.execute(statement)
	}
}

// VisitBlockStmt 处理代码块语句
func (i *Interpreter) VisitBlockStmt(stmt *ast.Block) interface{} {
	i.executeBlock(stmt.Statements, environment.NewEnclosedEnvironment(i.environment))
	return nil
}

// VisitExpressionStmt 处理表达式语句
func (i *Interpreter) VisitExpressionStmt(stmt *ast.Expression) interface{} {
	i.evaluate(stmt.Expr)
	return nil
}

// VisitIfStmt 处理条件语句
func (i *Interpreter) VisitIfStmt(stmt *ast.If) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

// VisitPrintStmt 处理打印语句
func (i *Interpreter) VisitPrintStmt(stmt *ast.Print) interface{} {
	value := i.evaluate(stmt.Expr)
	fmt.Println(i.stringify(value))
	return nil
}

// VisitVarStmt 处理变量声明语句
func (i *Interpreter) VisitVarStmt(stmt *ast.Var) interface{} {
	var value interface{}

	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

// VisitWhileStmt 处理循环语句
func (i *Interpreter) VisitWhileStmt(stmt *ast.While) interface{} {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

// handlePanic 处理解释过程中的异常
func (i *Interpreter) handlePanic() {
	if r := recover(); r != nil {
		if runtimeError, ok := r.(error.RuntimeError); ok {
			// 报告运行时错误
			i.errorReporter.Error(runtimeError.Token, 0, runtimeError.Message)
		} else {
			// 重新抛出其他异常
			panic(r)
		}
	}
}

// evaluate 求值表达式
func (i *Interpreter) evaluate(expr ast.Expr) interface{} {
	return expr.Accept(i)
}

// VisitAssignExpr 处理赋值表达式
func (i *Interpreter) VisitAssignExpr(expr *ast.Assign) interface{} {
	value := i.evaluate(expr.Value)
	i.environment.Assign(expr.Name, value)
	return value
}

// VisitLiteralExpr 处理字面量表达式
func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) interface{} {
	return expr.Value
}

// VisitGroupingExpr 处理分组表达式
func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

// VisitUnaryExpr 处理一元表达式
func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.MINUS:
		i.checkNumberOperand(expr.Operator, right)
		return -i.asNumber(right)
	case token.BANG:
		return !i.isTruthy(right)
	}

	// 不可达
	return nil
}

// VisitBinaryExpr 处理二元表达式
func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.MINUS:
		i.checkNumberOperands(expr.Operator, left, right)
		return i.asNumber(left) - i.asNumber(right)
	case token.SLASH:
		i.checkNumberOperands(expr.Operator, left, right)
		rightNum := i.asNumber(right)
		if rightNum == 0 {
			panic(error.RuntimeError{Token: expr.Operator, Message: "除数不能为零。"})
		}
		return i.asNumber(left) / rightNum
	case token.STAR:
		i.checkNumberOperands(expr.Operator, left, right)
		return i.asNumber(left) * i.asNumber(right)
	case token.PLUS:
		if i.isNumber(left) && i.isNumber(right) {
			return i.asNumber(left) + i.asNumber(right)
		}
		if i.isString(left) && i.isString(right) {
			return i.asString(left) + i.asString(right)
		}
		panic(error.RuntimeError{Token: expr.Operator, Message: "'+'运算符只能用于数字或字符串。"})
	case token.GREATER:
		i.checkNumberOperands(expr.Operator, left, right)
		return i.asNumber(left) > i.asNumber(right)
	case token.GREATER_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return i.asNumber(left) >= i.asNumber(right)
	case token.LESS:
		i.checkNumberOperands(expr.Operator, left, right)
		return i.asNumber(left) < i.asNumber(right)
	case token.LESS_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return i.asNumber(left) <= i.asNumber(right)
	case token.BANG_EQUAL:
		return !i.isEqual(left, right)
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	case token.COMMA:
		// 逗号表达式，返回右侧值
		return right
	}

	// 不可达
	return nil
}

// VisitTernaryExpr 处理三元表达式
func (i *Interpreter) VisitTernaryExpr(expr *ast.Ternary) interface{} {
	condition := i.evaluate(expr.Condition)

	if i.isTruthy(condition) {
		return i.evaluate(expr.ThenBranch)
	} else {
		return i.evaluate(expr.ElseBranch)
	}
}

// VisitVariableExpr 处理变量表达式
func (i *Interpreter) VisitVariableExpr(expr *ast.Variable) interface{} {
	return i.environment.Get(expr.Name)
}

// VisitLogicalExpr 处理逻辑表达式
func (i *Interpreter) VisitLogicalExpr(expr *ast.Logical) interface{} {
	left := i.evaluate(expr.Left)

	// 逻辑运算的短路处理
	if expr.Operator.Type == token.OR {
		if i.isTruthy(left) {
			return left // 短路求值
		}
	} else { // AND
		if !i.isTruthy(left) {
			return left // 短路求值
		}
	}

	// 如果没有短路，计算右侧的值
	return i.evaluate(expr.Right)
}

// 工具方法

// isTruthy 判断一个值是否为真
func (i *Interpreter) isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return true
}

// isEqual 判断两个值是否相等
func (i *Interpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

// isNumber 判断一个值是否为数字
func (i *Interpreter) isNumber(value interface{}) bool {
	_, ok := value.(float64)
	return ok
}

// isString 判断一个值是否为字符串
func (i *Interpreter) isString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

// asNumber 转换为数字
func (i *Interpreter) asNumber(value interface{}) float64 {
	number, ok := value.(float64)
	if !ok {
		// 应该不会到这里，因为我们已经用checkNumberOperand检查过了
		return 0
	}
	return number
}

// asString 转换为字符串
func (i *Interpreter) asString(value interface{}) string {
	str, ok := value.(string)
	if !ok {
		return ""
	}
	return str
}

// checkNumberOperand 检查一元运算符的操作数是否为数字
func (i *Interpreter) checkNumberOperand(operator *token.Token, operand interface{}) {
	if i.isNumber(operand) {
		return
	}
	panic(error.RuntimeError{Token: operator, Message: "操作数必须是数字。"})
}

// checkNumberOperands 检查二元运算符的操作数是否为数字
func (i *Interpreter) checkNumberOperands(operator *token.Token, left, right interface{}) {
	if i.isNumber(left) && i.isNumber(right) {
		return
	}
	panic(error.RuntimeError{Token: operator, Message: "操作数必须是数字。"})
}

// stringify 将值转换为字符串
func (i *Interpreter) stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}

	if num, ok := value.(float64); ok {
		text := fmt.Sprintf("%g", num)
		// 移除整数的小数点
		if math.Floor(num) == num {
			text = fmt.Sprintf("%.0f", num)
		}
		return text
	}

	return fmt.Sprintf("%v", value)
}
