package ast

import (
	"github.com/aixiasang/goLox/lox/token"
)

// Expr 是所有表达式节点的接口
type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

// ExprVisitor 定义了访问者模式的接口
type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
}

// Binary 表示二元表达式
type Binary struct {
	Left     Expr         // 左操作数
	Operator *token.Token // 操作符
	Right    Expr         // 右操作数
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

// Grouping 表示括号分组的表达式
type Grouping struct {
	Expression Expr // 括号内的表达式
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

// Literal 表示字面量
type Literal struct {
	Value interface{} // 字面量的值
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

// Unary 表示一元表达式
type Unary struct {
	Operator *token.Token // 操作符
	Right    Expr         // 操作数
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

// 创建表达式节点的构造函数
func NewBinary(left Expr, operator *token.Token, right Expr) *Binary {
	return &Binary{Left: left, Operator: operator, Right: right}
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{Expression: expression}
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{Value: value}
}

func NewUnary(operator *token.Token, right Expr) *Unary {
	return &Unary{Operator: operator, Right: right}
}
