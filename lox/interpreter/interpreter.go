package interpreter

import (
	"fmt"
	"math"

	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/environment"
	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/token"
)

// BreakException 表示遇到break语句的异常
type BreakException struct{}

// Interpreter 实现表达式求值和语句执行
type Interpreter struct {
	errorReporter error.Reporter
	environment   *environment.Environment
	locals        map[ast.Expr]int         // 变量的作用域深度信息
	globals       *environment.Environment // 全局环境
}

// NewInterpreter 创建一个新的解释器
func NewInterpreter(errorReporter error.Reporter) *Interpreter {
	globals := environment.NewEnvironment()

	// 添加内置函数
	globals.Define("clock", &Clock{})

	interpreter := &Interpreter{
		errorReporter: errorReporter,
		environment:   globals,
		globals:       globals,
		locals:        make(map[ast.Expr]int),
	}

	return interpreter
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
	defer func() {
		// 捕获break异常，但不做处理，仅用于跳出循环
		if r := recover(); r != nil {
			// 如果是BreakException，就直接吞掉
			if _, ok := r.(BreakException); !ok {
				// 其他类型的异常继续向上抛
				panic(r)
			}
		}
	}()

	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

// VisitBreakStmt 处理break语句
func (i *Interpreter) VisitBreakStmt(stmt *ast.Break) interface{} {
	// 通过抛出异常来跳出循环
	panic(BreakException{})
}

// VisitFunctionStmt 处理函数声明语句
func (i *Interpreter) VisitFunctionStmt(stmt *ast.Function) interface{} {
	function := NewFunction(stmt, i.environment)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil
}

// VisitReturnStmt 处理return语句
func (i *Interpreter) VisitReturnStmt(stmt *ast.Return) interface{} {
	var value interface{} = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	// 通过抛出ReturnValue异常来实现return语句
	panic(ReturnValue{Value: value})
}

// handlePanic 处理解释过程中的异常
func (i *Interpreter) handlePanic() {
	if r := recover(); r != nil {
		if runtimeError, ok := r.(error.RuntimeError); ok {
			// 报告运行时错误
			i.errorReporter.Error(runtimeError.Token, 0, runtimeError.Message)
		} else if _, ok := r.(BreakException); ok {
			// break语句超出循环范围
			// 这里不应该发生，因为解析器应该检查break是否在循环内
			i.errorReporter.ReportError(0, "Break语句只能在循环内部使用。")
		} else if _, ok := r.(ReturnValue); ok {
			// 返回值流动到最顶层
			// 这里可以选择将值作为REPL的结果返回
			return
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

	if distance, ok := i.locals[expr]; ok {
		i.environment.AssignAt(distance, expr.Name, value)
	} else {
		i.globals.Assign(expr.Name, value)
	}

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
	case token.MODULO:
		i.checkNumberOperands(expr.Operator, left, right)
		rightNum := i.asNumber(right)
		if rightNum == 0 {
			panic(error.RuntimeError{Token: expr.Operator, Message: "取模运算符的右操作数不能为零。"})
		}
		return float64(int(i.asNumber(left)) % int(i.asNumber(right)))
	case token.PLUS:
		if i.isNumber(left) && i.isNumber(right) {
			return i.asNumber(left) + i.asNumber(right)
		}
		// 如果任一操作数是字符串，则将另一个操作数也转换为字符串
		if i.isString(left) || i.isString(right) {
			return i.stringify(left) + i.stringify(right)
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
	return i.lookUpVariable(expr.Name, expr)
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

// VisitCallExpr 处理函数调用表达式
func (i *Interpreter) VisitCallExpr(expr *ast.Call) interface{} {
	callee := i.evaluate(expr.Callee)

	// 收集参数
	var arguments []interface{}
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	// 检查callee是否可调用
	function, ok := callee.(Callable)
	if !ok {
		panic(error.RuntimeError{Token: expr.Paren, Message: "只能调用函数和类。"})
	}

	// 检查参数数量是否正确
	if len(arguments) != function.Arity() {
		message := fmt.Sprintf("期望%d个参数，但得到%d个。", function.Arity(), len(arguments))
		panic(error.RuntimeError{Token: expr.Paren, Message: message})
	}

	// 调用函数
	return function.Call(i, arguments)
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

	// 如果是浮点数
	if num, ok := value.(float64); ok {
		text := fmt.Sprintf("%g", num)
		// 如果是整数，去掉小数点和小数点后的零
		if math.Floor(num) == num {
			text = fmt.Sprintf("%.0f", num)
		}
		return text
	}

	// 如果是字符串，直接返回
	if str, ok := value.(string); ok {
		return str
	}

	// 如果是布尔值
	if b, ok := value.(bool); ok {
		if b {
			return "true"
		}
		return "false"
	}

	// 如果是可调用对象，调用其String方法
	if callable, ok := value.(Callable); ok {
		return callable.String()
	}

	// 其他类型
	return fmt.Sprintf("%v", value)
}

// Resolve 记录变量引用的作用域深度
func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

// lookUpVariable 根据作用域深度查找变量
func (i *Interpreter) lookUpVariable(name *token.Token, expr ast.Expr) interface{} {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.GetAt(distance, name.Lexeme)
	} else {
		return i.globals.Get(name)
	}
}
