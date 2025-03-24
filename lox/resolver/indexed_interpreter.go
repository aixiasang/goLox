package resolver

import (
	"fmt"
	"math"
	"strconv"

	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/token"
)

// ReturnValue 表示return语句的返回值
type ReturnValue struct {
	Value interface{}
}

// BreakException 表示遇到break语句的异常
type BreakException struct{}

// IndexedEnvironment 是一个基于数组的环境，用于高效地访问变量
type IndexedEnvironment struct {
	values    []interface{}       // 当前作用域中的值
	enclosing *IndexedEnvironment // 外层环境
}

// NewIndexedEnvironment 创建一个新的索引环境
func NewIndexedEnvironment(enclosing *IndexedEnvironment) *IndexedEnvironment {
	return &IndexedEnvironment{
		values:    make([]interface{}, 0),
		enclosing: enclosing,
	}
}

// Define 定义一个新变量
func (e *IndexedEnvironment) Define(value interface{}) int {
	e.values = append(e.values, value)
	return len(e.values) - 1
}

// GetAt 获取指定深度和索引处的变量值
func (e *IndexedEnvironment) GetAt(depth int, index int) interface{} {
	if depth == 0 {
		return e.values[index]
	}

	// 递归地在外层环境中查找变量
	return e.ancestor(depth).values[index]
}

// AssignAt 在指定深度和索引处设置变量值
func (e *IndexedEnvironment) AssignAt(depth int, index int, value interface{}) {
	if depth == 0 {
		e.values[index] = value
		return
	}

	// 递归地在外层环境中设置变量
	e.ancestor(depth).values[index] = value
}

// 获取指定深度的环境
func (e *IndexedEnvironment) ancestor(depth int) *IndexedEnvironment {
	environment := e
	for i := 0; i < depth; i++ {
		environment = environment.enclosing
	}
	return environment
}

// IndexedInterpreter 使用索引访问变量的解释器
type IndexedInterpreter struct {
	errorReporter error.Reporter
	environment   *IndexedEnvironment
	globals       *IndexedEnvironment
	locals        map[ast.Expr]VarLocation
}

// NewIndexedInterpreter 创建一个新的索引优化解释器
func NewIndexedInterpreter(errorReporter error.Reporter) *IndexedInterpreter {
	globals := NewIndexedEnvironment(nil)

	// 添加内置函数和变量
	// 注意: 这种情况下，全局变量和函数仍然需要使用名称查找

	interpreter := &IndexedInterpreter{
		errorReporter: errorReporter,
		environment:   globals,
		globals:       globals,
		locals:        make(map[ast.Expr]VarLocation),
	}

	return interpreter
}

// SetLocations 设置变量位置信息
func (i *IndexedInterpreter) SetLocations(locals map[ast.Expr]VarLocation) {
	i.locals = locals
}

// Interpret 解释执行语句列表
func (i *IndexedInterpreter) Interpret(statements []ast.Stmt) {
	defer i.handlePanic()

	for _, stmt := range statements {
		i.execute(stmt)
	}
}

// execute 执行一条语句
func (i *IndexedInterpreter) execute(stmt ast.Stmt) {
	stmt.Accept(i)
}

// executeBlock 在给定环境中执行语句块
func (i *IndexedInterpreter) executeBlock(statements []ast.Stmt, env *IndexedEnvironment) {
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

// handlePanic 处理解释过程中的异常
func (i *IndexedInterpreter) handlePanic() {
	if r := recover(); r != nil {
		if runtimeError, ok := r.(error.RuntimeError); ok {
			// 报告运行时错误
			i.errorReporter.Error(runtimeError.Token, 0, runtimeError.Message)
		} else if _, ok := r.(BreakException); ok {
			// break语句超出循环范围
			i.errorReporter.ReportError(0, "Break语句只能在循环内部使用。")
		} else if _, ok := r.(ReturnValue); ok {
			// 返回值流动到最顶层
			return
		} else {
			// 重新抛出其他异常
			panic(r)
		}
	}
}

// lookUpVariable 查找变量的值
func (i *IndexedInterpreter) lookUpVariable(name *token.Token, expr ast.Expr) interface{} {
	if location, ok := i.locals[expr]; ok {
		return i.environment.GetAt(location.Depth, location.Index)
	}

	// 如果不是局部变量，可能是全局变量，但当前实现中没有全局变量表
	// 所以返回错误
	i.errorReporter.Error(name, 0, "未定义的变量。")
	return nil
}

// VisitBlockStmt 处理代码块语句
func (i *IndexedInterpreter) VisitBlockStmt(stmt *ast.Block) interface{} {
	i.executeBlock(stmt.Statements, NewEnclosedIndexedEnvironment(i.environment))
	return nil
}

// VisitExpressionStmt 处理表达式语句
func (i *IndexedInterpreter) VisitExpressionStmt(stmt *ast.Expression) interface{} {
	i.evaluate(stmt.Expr)
	return nil
}

// VisitIfStmt 处理条件语句
func (i *IndexedInterpreter) VisitIfStmt(stmt *ast.If) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

// VisitPrintStmt 处理打印语句
func (i *IndexedInterpreter) VisitPrintStmt(stmt *ast.Print) interface{} {
	value := i.evaluate(stmt.Expr)
	fmt.Println(i.stringify(value))
	return nil
}

// VisitVarStmt 处理变量声明语句
func (i *IndexedInterpreter) VisitVarStmt(stmt *ast.Var) interface{} {
	var value interface{}

	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	// 变量定义返回索引
	i.environment.Define(value)
	return nil
}

// VisitWhileStmt 处理循环语句
func (i *IndexedInterpreter) VisitWhileStmt(stmt *ast.While) interface{} {
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
func (i *IndexedInterpreter) VisitBreakStmt(stmt *ast.Break) interface{} {
	// 通过抛出异常来跳出循环
	panic(BreakException{})
}

// VisitFunctionStmt 处理函数声明语句
func (i *IndexedInterpreter) VisitFunctionStmt(stmt *ast.Function) interface{} {
	function := &LoxFunction{
		declaration: stmt,
		closure:     i.environment,
		interpreter: i,
	}

	i.environment.Define(function)
	return nil
}

// VisitReturnStmt 处理return语句
func (i *IndexedInterpreter) VisitReturnStmt(stmt *ast.Return) interface{} {
	var value interface{} = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	// 通过抛出ReturnValue异常来实现return语句
	panic(ReturnValue{Value: value})
}

// evaluate 求值表达式
func (i *IndexedInterpreter) evaluate(expr ast.Expr) interface{} {
	return expr.Accept(i)
}

// VisitAssignExpr 处理赋值表达式
func (i *IndexedInterpreter) VisitAssignExpr(expr *ast.Assign) interface{} {
	value := i.evaluate(expr.Value)

	if location, ok := i.locals[expr]; ok {
		i.environment.AssignAt(location.Depth, location.Index, value)
	} else {
		// 全局变量处理逻辑
	}

	return value
}

// VisitVariableExpr 处理变量表达式
func (i *IndexedInterpreter) VisitVariableExpr(expr *ast.Variable) interface{} {
	return i.lookUpVariable(expr.Name, expr)
}

// VisitLiteralExpr 处理字面量表达式
func (i *IndexedInterpreter) VisitLiteralExpr(expr *ast.Literal) interface{} {
	return expr.Value
}

// VisitGroupingExpr 处理分组表达式
func (i *IndexedInterpreter) VisitGroupingExpr(expr *ast.Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

// VisitUnaryExpr 处理一元表达式
func (i *IndexedInterpreter) VisitUnaryExpr(expr *ast.Unary) interface{} {
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
func (i *IndexedInterpreter) VisitBinaryExpr(expr *ast.Binary) interface{} {
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
	}

	// 不可达
	return nil
}

// VisitLogicalExpr 处理逻辑表达式
func (i *IndexedInterpreter) VisitLogicalExpr(expr *ast.Logical) interface{} {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == token.OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

// VisitCallExpr 处理函数调用表达式
func (i *IndexedInterpreter) VisitCallExpr(expr *ast.Call) interface{} {
	callee := i.evaluate(expr.Callee)

	arguments := make([]interface{}, len(expr.Arguments))
	for j, argument := range expr.Arguments {
		arguments[j] = i.evaluate(argument)
	}

	if function, ok := callee.(LoxCallable); ok {
		if len(arguments) != function.Arity() {
			i.errorReporter.Error(expr.Paren, 0,
				fmt.Sprintf("期望 %d 个参数但得到 %d 个。", function.Arity(), len(arguments)))
			return nil
		}

		return function.Call(i, arguments)
	}

	i.errorReporter.Error(expr.Paren, 0, "只能调用函数和类。")
	return nil
}

// VisitTernaryExpr 处理三元表达式
func (i *IndexedInterpreter) VisitTernaryExpr(expr *ast.Ternary) interface{} {
	condition := i.evaluate(expr.Condition)

	if i.isTruthy(condition) {
		return i.evaluate(expr.ThenBranch)
	} else {
		return i.evaluate(expr.ElseBranch)
	}
}

// 辅助方法

// isTruthy 确定值是否为"真"
func (i *IndexedInterpreter) isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return true
}

// isEqual 比较两个值是否相等
func (i *IndexedInterpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

// checkNumberOperand 检查操作数是否为数字
func (i *IndexedInterpreter) checkNumberOperand(operator *token.Token, operand interface{}) {
	if i.isNumber(operand) {
		return
	}
	panic(error.RuntimeError{Token: operator, Message: "操作数必须是数字。"})
}

// checkNumberOperands 检查操作数是否为数字
func (i *IndexedInterpreter) checkNumberOperands(operator *token.Token, left, right interface{}) {
	if i.isNumber(left) && i.isNumber(right) {
		return
	}
	panic(error.RuntimeError{Token: operator, Message: "操作数必须是数字。"})
}

// isNumber 检查值是否为数字
func (i *IndexedInterpreter) isNumber(value interface{}) bool {
	_, ok := value.(float64)
	return ok
}

// asNumber 将值转换为数字
func (i *IndexedInterpreter) asNumber(value interface{}) float64 {
	if num, ok := value.(float64); ok {
		return num
	}
	if str, ok := value.(string); ok {
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			return num
		}
	}
	return 0
}

// isString 检查值是否为字符串
func (i *IndexedInterpreter) isString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

// stringify 将值转换为字符串
func (i *IndexedInterpreter) stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}

	if num, ok := value.(float64); ok {
		text := fmt.Sprintf("%g", num)
		if math.Floor(num) == num {
			text = fmt.Sprintf("%.0f", num)
		}
		return text
	}

	return fmt.Sprintf("%v", value)
}

// LoxCallable 接口定义了可调用的Lox对象
type LoxCallable interface {
	Call(interpreter *IndexedInterpreter, arguments []interface{}) interface{}
	Arity() int
}

// LoxFunction 表示Lox函数
type LoxFunction struct {
	declaration   *ast.Function
	closure       *IndexedEnvironment
	interpreter   *IndexedInterpreter
	isInitializer bool
}

// Arity 返回函数参数数量
func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

// Call 调用函数
func (f *LoxFunction) Call(interpreter *IndexedInterpreter, arguments []interface{}) interface{} {
	environment := NewIndexedEnvironment(f.closure)

	// 为函数的每个参数定义变量
	for i := 0; i < len(f.declaration.Params); i++ {
		environment.Define(arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(ReturnValue); ok {
				if f.isInitializer {
					// 如果是初始化方法，始终返回this
					return
				}
				panic(returnValue)
			}
			// 重新抛出非返回值的panic
			panic(r)
		}
	}()

	interpreter.executeBlock(f.declaration.Body, environment)

	if f.isInitializer {
		// 返回this
		return f.closure.GetAt(0, 0)
	}

	return nil
}

// Bind 将方法绑定到实例
func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := NewIndexedEnvironment(f.closure)
	environment.Define(instance)

	return &LoxFunction{
		declaration:   f.declaration,
		closure:       environment,
		interpreter:   f.interpreter,
		isInitializer: f.isInitializer,
	}
}

// LoxClass 表示Lox类
type LoxClass struct {
	name        string
	superclass  *LoxClass
	methods     map[string]*LoxFunction
	interpreter *IndexedInterpreter
}

// Arity 返回类构造函数参数数量
func (c *LoxClass) Arity() int {
	initializer := c.FindMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

// Call 调用类构造函数
func (c *LoxClass) Call(interpreter *IndexedInterpreter, arguments []interface{}) interface{} {
	instance := &LoxInstance{
		class:  c,
		fields: make(map[string]interface{}),
	}

	initializer := c.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}

	return instance
}

// FindMethod 查找类方法
func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if method, ok := c.methods[name]; ok {
		return method
	}

	if c.superclass != nil {
		return c.superclass.FindMethod(name)
	}

	return nil
}

// LoxInstance 表示Lox类实例
type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

// Get 获取实例属性
func (i *LoxInstance) Get(name *token.Token) interface{} {
	if field, ok := i.fields[name.Lexeme]; ok {
		return field
	}

	if method := i.class.FindMethod(name.Lexeme); method != nil {
		return method.Bind(i)
	}

	i.class.interpreter.errorReporter.Error(name, 0,
		fmt.Sprintf("未定义属性 '%s'。", name.Lexeme))
	return nil
}

// Set 设置实例属性
func (i *LoxInstance) Set(name *token.Token, value interface{}) {
	i.fields[name.Lexeme] = value
}
