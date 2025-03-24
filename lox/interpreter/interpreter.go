package interpreter

import (
	"fmt"
	"math"

	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/token"
)

// Interpreter 实现表达式求值
type Interpreter struct {
	errorReporter error.Reporter
}

// NewInterpreter 创建一个新的解释器
func NewInterpreter(errorReporter error.Reporter) *Interpreter {
	return &Interpreter{
		errorReporter: errorReporter,
	}
}

// Interpret 解释执行表达式
func (i *Interpreter) Interpret(expr ast.Expr) {
	defer i.handlePanic()
	value := i.evaluate(expr)
	fmt.Println(i.stringify(value))
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
	// 目前我们还没有环境来存储变量，返回nil
	panic(error.RuntimeError{Token: expr.Name, Message: "未定义的变量。"})
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
