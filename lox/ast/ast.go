package ast

import (
	"github.com/aixiasang/goLox/lox/token"
)

// Expr 表达式接口
type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

// ExprVisitor 访问者接口
type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitTernaryExpr(expr *Ternary) interface{}
	VisitVariableExpr(expr *Variable) interface{}
	VisitAssignExpr(expr *Assign) interface{}
	VisitLogicalExpr(expr *Logical) interface{}
}

// Binary 二元表达式
type Binary struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

// Accept 接受访问者
func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

// NewBinary 创建二元表达式
func NewBinary(left Expr, operator *token.Token, right Expr) *Binary {
	return &Binary{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// Grouping 分组表达式
type Grouping struct {
	Expression Expr
}

// Accept 接受访问者
func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

// NewGrouping 创建分组表达式
func NewGrouping(expression Expr) *Grouping {
	return &Grouping{
		Expression: expression,
	}
}

// Literal 字面量表达式
type Literal struct {
	Value interface{}
}

// Accept 接受访问者
func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

// NewLiteral 创建字面量表达式
func NewLiteral(value interface{}) *Literal {
	return &Literal{
		Value: value,
	}
}

// Unary 一元表达式
type Unary struct {
	Operator *token.Token
	Right    Expr
}

// Accept 接受访问者
func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

// NewUnary 创建一元表达式
func NewUnary(operator *token.Token, right Expr) *Unary {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

// Ternary 三元条件表达式
type Ternary struct {
	Condition  Expr
	ThenBranch Expr
	ElseBranch Expr
}

// Accept 接受访问者
func (t *Ternary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitTernaryExpr(t)
}

// NewTernary 创建三元表达式
func NewTernary(condition Expr, thenBranch Expr, elseBranch Expr) *Ternary {
	return &Ternary{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

// Variable 变量表达式
type Variable struct {
	Name *token.Token
}

// Accept 接受访问者
func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}

// NewVariable 创建变量表达式
func NewVariable(name *token.Token) *Variable {
	return &Variable{
		Name: name,
	}
}

// Assign 赋值表达式
type Assign struct {
	Name  *token.Token
	Value Expr
}

// Accept 接受访问者
func (a *Assign) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(a)
}

// NewAssign 创建赋值表达式
func NewAssign(name *token.Token, value Expr) *Assign {
	return &Assign{
		Name:  name,
		Value: value,
	}
}

// Logical 逻辑表达式
type Logical struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

// Accept 接受访问者
func (l *Logical) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogicalExpr(l)
}

// NewLogical 创建逻辑表达式
func NewLogical(left Expr, operator *token.Token, right Expr) *Logical {
	return &Logical{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}
