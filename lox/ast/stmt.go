package ast

import (
	"github.com/aixiasang/goLox/lox/token"
)

// Stmt 语句接口
type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

// StmtVisitor 语句访问者接口
type StmtVisitor interface {
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitPrintStmt(stmt *Print) interface{}
	VisitVarStmt(stmt *Var) interface{}
	VisitBlockStmt(stmt *Block) interface{}
	VisitIfStmt(stmt *If) interface{}
	VisitWhileStmt(stmt *While) interface{}
}

// Expression 表达式语句
type Expression struct {
	Expr Expr
}

// Accept 接受访问者
func (e *Expression) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(e)
}

// NewExpression 创建表达式语句
func NewExpression(expr Expr) *Expression {
	return &Expression{
		Expr: expr,
	}
}

// Print 打印语句
type Print struct {
	Expr Expr
}

// Accept 接受访问者
func (p *Print) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(p)
}

// NewPrint 创建打印语句
func NewPrint(expr Expr) *Print {
	return &Print{
		Expr: expr,
	}
}

// Var 变量声明语句
type Var struct {
	Name        *token.Token
	Initializer Expr
}

// Accept 接受访问者
func (v *Var) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(v)
}

// NewVar 创建变量声明语句
func NewVar(name *token.Token, initializer Expr) *Var {
	return &Var{
		Name:        name,
		Initializer: initializer,
	}
}

// Block 代码块语句
type Block struct {
	Statements []Stmt
}

// Accept 接受访问者
func (b *Block) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitBlockStmt(b)
}

// NewBlock 创建代码块语句
func NewBlock(statements []Stmt) *Block {
	return &Block{
		Statements: statements,
	}
}

// If 条件语句
type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt // 可能为nil
}

// Accept 接受访问者
func (i *If) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitIfStmt(i)
}

// NewIf 创建条件语句
func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt) *If {
	return &If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

// While 循环语句
type While struct {
	Condition Expr
	Body      Stmt
}

// Accept 接受访问者
func (w *While) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitWhileStmt(w)
}

// NewWhile 创建循环语句
func NewWhile(condition Expr, body Stmt) *While {
	return &While{
		Condition: condition,
		Body:      body,
	}
}
