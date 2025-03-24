package resolver

import (
	"fmt"
	"strings"

	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/interpreter"
	"github.com/aixiasang/goLox/lox/token"
)

// FunctionType 函数类型枚举
// type FunctionType int

// const (
// 	// FunctionNONE 非函数上下文
// 	FunctionNONE FunctionType = iota
// 	// FunctionFUNCTION 普通函数上下文
// 	FunctionFUNCTION
// )

// Resolver 静态分析和变量解析器
type Resolver struct {
	interpreter     *interpreter.Interpreter
	errorReporter   error.Reporter
	scopes          []map[string]bool // 作用域栈
	currentFunction FunctionType      // 当前函数上下文
	locals          map[string]bool   // 追踪变量是否被使用
}

// NewResolver 创建一个新的Resolver
func NewResolver(interpreter *interpreter.Interpreter, errorReporter error.Reporter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		errorReporter:   errorReporter,
		scopes:          make([]map[string]bool, 0),
		currentFunction: FunctionNONE,
		locals:          make(map[string]bool),
	}
}

// Resolve 解析一组语句
func (r *Resolver) Resolve(statements []ast.Stmt) {
	for _, stmt := range statements {
		r.resolveStmt(stmt)
	}
}

// resolveStmt 解析单个语句
func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

// resolveExpr 解析表达式
func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

// beginScope 开始一个新的作用域
func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

// endScope 结束当前作用域
func (r *Resolver) endScope() {
	if len(r.scopes) > 0 {
		// 检查未使用的变量
		scope := r.scopes[len(r.scopes)-1]
		scopeDepth := len(r.scopes) - 1

		// 遍历当前作用域中的所有变量
		for name := range scope {
			varKey := fmt.Sprintf("%s:%d", name, scopeDepth)
			// 如果变量已声明但从未使用过
			if used, exists := r.locals[varKey]; exists && !used {
				// 忽略函数名和循环变量（这是策略选择，可以根据需要调整）
				if !strings.HasPrefix(name, "fn_") && !strings.HasPrefix(name, "loop_") {
					r.errorReporter.ReportError(0, fmt.Sprintf("局部变量 '%s' 已声明但从未使用", name))
				}
			}
		}

		r.scopes = r.scopes[:len(r.scopes)-1]
	}
}

// declare 声明一个变量
func (r *Resolver) declare(name *token.Token) {
	if len(r.scopes) == 0 {
		return // 全局作用域
	}

	scope := r.scopes[len(r.scopes)-1]

	// 检查变量是否已经在当前作用域中声明
	if _, exists := scope[name.Lexeme]; exists {
		r.errorReporter.Error(name, 0, "变量 '"+name.Lexeme+"' 已在此作用域中声明。")
	}

	// 标记为"尚未初始化"
	scope[name.Lexeme] = false

	// 标记变量为未使用
	scopeDepth := len(r.scopes) - 1
	varKey := fmt.Sprintf("%s:%d", name.Lexeme, scopeDepth)
	r.locals[varKey] = false
}

// define 定义一个变量（完成初始化）
func (r *Resolver) define(name *token.Token) {
	if len(r.scopes) == 0 {
		return // 全局作用域
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true // 标记为"已初始化"
}

// resolveLocal 解析局部变量
func (r *Resolver) resolveLocal(expr ast.Expr, name *token.Token) {
	// 从内向外查找变量声明
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, exists := r.scopes[i][name.Lexeme]; exists {
			// 找到变量，计算深度
			depth := len(r.scopes) - 1 - i
			r.interpreter.Resolve(expr, depth)

			// 标记变量为已使用
			varKey := fmt.Sprintf("%s:%d", name.Lexeme, i)
			r.locals[varKey] = true

			return
		}
	}

	// 如果这里没有找到，假设它是一个全局变量
}

// resolveFunction 解析函数声明
func (r *Resolver) resolveFunction(function *ast.Function, funcType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = funcType

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.Resolve(function.Body)
	r.endScope()

	r.currentFunction = enclosingFunction
}

// VisitBlockStmt 访问代码块
func (r *Resolver) VisitBlockStmt(stmt *ast.Block) interface{} {
	r.beginScope()
	r.Resolve(stmt.Statements)
	r.endScope()
	return nil
}

// VisitVarStmt 访问变量声明
func (r *Resolver) VisitVarStmt(stmt *ast.Var) interface{} {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

// VisitVariableExpr 访问变量引用
func (r *Resolver) VisitVariableExpr(expr *ast.Variable) interface{} {
	if len(r.scopes) > 0 {
		scope := r.scopes[len(r.scopes)-1]
		if initialized, exists := scope[expr.Name.Lexeme]; exists && !initialized {
			r.errorReporter.Error(expr.Name, 0, "不能在变量初始化中引用自身。")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

// VisitAssignExpr 访问赋值表达式
func (r *Resolver) VisitAssignExpr(expr *ast.Assign) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

// VisitFunctionStmt 访问函数声明
func (r *Resolver) VisitFunctionStmt(stmt *ast.Function) interface{} {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FunctionFUNCTION)
	return nil
}

// VisitExpressionStmt 访问表达式语句
func (r *Resolver) VisitExpressionStmt(stmt *ast.Expression) interface{} {
	r.resolveExpr(stmt.Expr)
	return nil
}

// VisitIfStmt 访问if语句
func (r *Resolver) VisitIfStmt(stmt *ast.If) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

// VisitPrintStmt 访问打印语句
func (r *Resolver) VisitPrintStmt(stmt *ast.Print) interface{} {
	r.resolveExpr(stmt.Expr)
	return nil
}

// VisitReturnStmt 访问return语句
func (r *Resolver) VisitReturnStmt(stmt *ast.Return) interface{} {
	if r.currentFunction == FunctionNONE {
		r.errorReporter.Error(stmt.Keyword, 0, "不能在函数外部使用return语句。")
	}

	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}

	return nil
}

// VisitWhileStmt 访问while语句
func (r *Resolver) VisitWhileStmt(stmt *ast.While) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

// VisitBreakStmt 访问break语句
func (r *Resolver) VisitBreakStmt(stmt *ast.Break) interface{} {
	return nil
}

// VisitBinaryExpr 访问二元表达式
func (r *Resolver) VisitBinaryExpr(expr *ast.Binary) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

// VisitCallExpr 访问函数调用表达式
func (r *Resolver) VisitCallExpr(expr *ast.Call) interface{} {
	r.resolveExpr(expr.Callee)

	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}

	return nil
}

// VisitGroupingExpr 访问分组表达式
func (r *Resolver) VisitGroupingExpr(expr *ast.Grouping) interface{} {
	r.resolveExpr(expr.Expression)
	return nil
}

// VisitLiteralExpr 访问字面量表达式
func (r *Resolver) VisitLiteralExpr(expr *ast.Literal) interface{} {
	return nil
}

// VisitLogicalExpr 访问逻辑表达式
func (r *Resolver) VisitLogicalExpr(expr *ast.Logical) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

// VisitUnaryExpr 访问一元表达式
func (r *Resolver) VisitUnaryExpr(expr *ast.Unary) interface{} {
	r.resolveExpr(expr.Right)
	return nil
}

// VisitTernaryExpr 访问三元表达式
func (r *Resolver) VisitTernaryExpr(expr *ast.Ternary) interface{} {
	r.resolveExpr(expr.Condition)
	r.resolveExpr(expr.ThenBranch)
	r.resolveExpr(expr.ElseBranch)
	return nil
}
