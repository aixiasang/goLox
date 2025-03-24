package resolver

import (
	"fmt"

	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/token"
)

// IndexedEnvironment 基于数组的环境
// type IndexedEnvironment struct {
// 	enclosing *IndexedEnvironment
// 	values    []interface{}
// }

// NewIndexedEnvironment 创建一个新的索引环境
// func NewIndexedEnvironment() *IndexedEnvironment {
// 	return &IndexedEnvironment{
// 		enclosing: nil,
// 		values:    make([]interface{}, 0),
// 	}
// }

// NewEnclosedIndexedEnvironment 创建一个包含外围环境的新索引环境
func NewEnclosedIndexedEnvironment(enclosing *IndexedEnvironment) *IndexedEnvironment {
	env := NewIndexedEnvironment(enclosing)
	return env
}

// Define 定义变量，返回变量在当前环境中的索引
// func (e *IndexedEnvironment) Define(value interface{}) int {
// 	e.values = append(e.values, value)
// 	return len(e.values) - 1
// }

// Get 通过索引获取变量值
func (e *IndexedEnvironment) Get(index int) interface{} {
	if index < len(e.values) {
		return e.values[index]
	}
	return nil
}

// GetAt 从指定深度的环境中通过索引获取变量值
// func (e *IndexedEnvironment) GetAt(distance int, index int) interface{} {
// 	return e.Ancestor(distance).Get(index)
// }

// Assign 通过索引给变量赋值
func (e *IndexedEnvironment) Assign(index int, value interface{}) {
	if index < len(e.values) {
		e.values[index] = value
	}
}

// AssignAt 在指定深度的环境中通过索引给变量赋值
// func (e *IndexedEnvironment) AssignAt(distance int, index int, value interface{}) {
// 	e.Ancestor(distance).Assign(index, value)
// }

// Ancestor 获取指定深度的环境
func (e *IndexedEnvironment) Ancestor(distance int) *IndexedEnvironment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

// FunctionType 函数类型枚举
type FunctionType int

const (
	// FunctionNONE 非函数上下文
	FunctionNONE FunctionType = iota
	// FunctionFUNCTION 普通函数上下文
	FunctionFUNCTION
)

// OptimizedResolver 优化的变量解析器
type OptimizedResolver struct {
	errorReporter   error.Reporter
	scopes          []map[string]VarInfo     // 作用域栈
	currentScope    int                      // 当前作用域级别
	variableCount   int                      // 当前作用域中的变量计数
	locations       map[ast.Expr]VarLocation // 变量位置信息
	currentFunction FunctionType             // 当前函数上下文
}

// VarInfo 变量信息
type VarInfo struct {
	Index       int  // 变量在作用域中的索引
	Initialized bool // 变量是否已初始化
	Used        bool // 变量是否被使用
}

// VarLocation 变量位置信息
type VarLocation struct {
	Depth int // 作用域深度
	Index int // 在该作用域中的索引
}

// NewOptimizedResolver 创建一个新的优化解析器
func NewOptimizedResolver(errorReporter error.Reporter) *OptimizedResolver {
	return &OptimizedResolver{
		errorReporter:   errorReporter,
		scopes:          make([]map[string]VarInfo, 0),
		currentScope:    -1,
		variableCount:   0,
		locations:       make(map[ast.Expr]VarLocation),
		currentFunction: FunctionNONE,
	}
}

// beginScope 开始一个新的作用域
func (r *OptimizedResolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]VarInfo))
	r.currentScope++
	r.variableCount = 0
}

// endScope 结束当前作用域
func (r *OptimizedResolver) endScope() {
	// 检查未使用的变量
	scope := r.scopes[r.currentScope]
	for name, info := range scope {
		if !info.Used {
			r.errorReporter.ReportError(0, fmt.Sprintf("局部变量 '%s' 已声明但从未使用", name))
		}
	}

	r.scopes = r.scopes[:len(r.scopes)-1]
	r.currentScope--
}

// declare 声明一个变量
func (r *OptimizedResolver) declare(name *token.Token) int {
	if r.currentScope < 0 {
		return -1 // 全局作用域
	}

	scope := r.scopes[r.currentScope]

	// 检查变量是否已经在当前作用域中声明
	if _, exists := scope[name.Lexeme]; exists {
		r.errorReporter.Error(name, 0, "变量 '"+name.Lexeme+"' 已在此作用域中声明。")
		return -1
	}

	// 分配一个新的索引
	index := r.variableCount
	r.variableCount++

	// 标记为"尚未初始化"且"未使用"
	scope[name.Lexeme] = VarInfo{
		Index:       index,
		Initialized: false,
		Used:        false,
	}

	return index
}

// define 定义一个变量（完成初始化）
func (r *OptimizedResolver) define(name *token.Token) {
	if r.currentScope < 0 {
		return // 全局作用域
	}

	scope := r.scopes[r.currentScope]
	if info, exists := scope[name.Lexeme]; exists {
		info.Initialized = true
		scope[name.Lexeme] = info
	}
}

// resolveLocal 解析局部变量，返回(深度，索引)
func (r *OptimizedResolver) resolveLocal(expr ast.Expr, name *token.Token) {
	// 从内向外查找变量声明
	for i := r.currentScope; i >= 0; i-- {
		scope := r.scopes[i]
		if info, exists := scope[name.Lexeme]; exists {
			// 检查变量是否在使用前已初始化
			if i == r.currentScope && !info.Initialized {
				r.errorReporter.Error(name, 0, "不能在变量初始化中引用自身。")
				return
			}

			// 标记变量为已使用
			info.Used = true
			scope[name.Lexeme] = info

			// 找到变量，计算深度和索引
			depth := r.currentScope - i
			r.locations[expr] = VarLocation{
				Depth: depth,
				Index: info.Index,
			}
			return
		}
	}

	// 如果这里没有找到，假设它是一个全局变量
}

// resolveFunction 解析函数声明
func (r *OptimizedResolver) resolveFunction(function *ast.Function, funcType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = funcType

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStatements(function.Body)
	r.endScope()

	r.currentFunction = enclosingFunction
}

// ResolveStatements 解析一组语句
func (r *OptimizedResolver) ResolveStatements(statements []ast.Stmt) map[ast.Expr]VarLocation {
	for _, stmt := range statements {
		r.resolveStmt(stmt)
	}
	return r.locations
}

// resolveStmt 解析单个语句
func (r *OptimizedResolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

// resolveExpr 解析表达式
func (r *OptimizedResolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

// VisitBlockStmt 访问代码块
func (r *OptimizedResolver) VisitBlockStmt(stmt *ast.Block) interface{} {
	r.beginScope()
	r.ResolveStatements(stmt.Statements)
	r.endScope()
	return nil
}

// VisitVarStmt 访问变量声明
func (r *OptimizedResolver) VisitVarStmt(stmt *ast.Var) interface{} {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

// VisitVariableExpr 访问变量引用
func (r *OptimizedResolver) VisitVariableExpr(expr *ast.Variable) interface{} {
	if r.currentScope >= 0 {
		scope := r.scopes[r.currentScope]
		if info, exists := scope[expr.Name.Lexeme]; exists && !info.Initialized {
			r.errorReporter.Error(expr.Name, 0, "不能在变量初始化中引用自身。")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

// VisitAssignExpr 访问赋值表达式
func (r *OptimizedResolver) VisitAssignExpr(expr *ast.Assign) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

// VisitFunctionStmt 访问函数声明
func (r *OptimizedResolver) VisitFunctionStmt(stmt *ast.Function) interface{} {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FunctionFUNCTION)
	return nil
}

// VisitExpressionStmt 访问表达式语句
func (r *OptimizedResolver) VisitExpressionStmt(stmt *ast.Expression) interface{} {
	r.resolveExpr(stmt.Expr)
	return nil
}

// VisitIfStmt 访问if语句
func (r *OptimizedResolver) VisitIfStmt(stmt *ast.If) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

// VisitPrintStmt 访问打印语句
func (r *OptimizedResolver) VisitPrintStmt(stmt *ast.Print) interface{} {
	r.resolveExpr(stmt.Expr)
	return nil
}

// VisitReturnStmt 访问return语句
func (r *OptimizedResolver) VisitReturnStmt(stmt *ast.Return) interface{} {
	if r.currentFunction == FunctionNONE {
		r.errorReporter.Error(stmt.Keyword, 0, "不能在函数外部使用return语句。")
	}

	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}

	return nil
}

// VisitWhileStmt 访问while语句
func (r *OptimizedResolver) VisitWhileStmt(stmt *ast.While) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

// VisitBreakStmt 访问break语句
func (r *OptimizedResolver) VisitBreakStmt(stmt *ast.Break) interface{} {
	return nil
}

// VisitBinaryExpr 访问二元表达式
func (r *OptimizedResolver) VisitBinaryExpr(expr *ast.Binary) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

// VisitCallExpr 访问函数调用表达式
func (r *OptimizedResolver) VisitCallExpr(expr *ast.Call) interface{} {
	r.resolveExpr(expr.Callee)

	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}

	return nil
}

// VisitGroupingExpr 访问分组表达式
func (r *OptimizedResolver) VisitGroupingExpr(expr *ast.Grouping) interface{} {
	r.resolveExpr(expr.Expression)
	return nil
}

// VisitLiteralExpr 访问字面量表达式
func (r *OptimizedResolver) VisitLiteralExpr(expr *ast.Literal) interface{} {
	return nil
}

// VisitLogicalExpr 访问逻辑表达式
func (r *OptimizedResolver) VisitLogicalExpr(expr *ast.Logical) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

// VisitUnaryExpr 访问一元表达式
func (r *OptimizedResolver) VisitUnaryExpr(expr *ast.Unary) interface{} {
	r.resolveExpr(expr.Right)
	return nil
}

// VisitTernaryExpr 访问三元表达式
func (r *OptimizedResolver) VisitTernaryExpr(expr *ast.Ternary) interface{} {
	r.resolveExpr(expr.Condition)
	r.resolveExpr(expr.ThenBranch)
	r.resolveExpr(expr.ElseBranch)
	return nil
}

// GetLocations 获取所有变量的位置信息
func (r *OptimizedResolver) GetLocations() map[ast.Expr]VarLocation {
	return r.locations
}
