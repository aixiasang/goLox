package ast

import (
	"fmt"
	"strings"
)

// RpnPrinter 实现了表达式访问者接口，用于将AST转换为逆波兰表示法(RPN)字符串
type RpnPrinter struct{}

// NewRpnPrinter 创建一个新的RPN打印器
func NewRpnPrinter() *RpnPrinter {
	return &RpnPrinter{}
}

// Print 将表达式转换为RPN字符串
func (p *RpnPrinter) Print(expr Expr) string {
	return expr.Accept(p).(string)
}

// VisitBinaryExpr 访问二元表达式
// 在RPN中，操作数在操作符之前
func (p *RpnPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	var builder strings.Builder

	// 先处理左操作数
	builder.WriteString(expr.Left.Accept(p).(string))
	builder.WriteString(" ")

	// 再处理右操作数
	builder.WriteString(expr.Right.Accept(p).(string))
	builder.WriteString(" ")

	// 最后是操作符
	builder.WriteString(expr.Operator.Lexeme)

	return builder.String()
}

// VisitGroupingExpr 访问分组表达式
// 在RPN中，括号不需要显式表示，因为操作顺序已经由位置决定
func (p *RpnPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return expr.Expression.Accept(p)
}

// VisitLiteralExpr 访问字面量
func (p *RpnPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

// VisitUnaryExpr 访问一元表达式
// 在RPN中，一元操作符在操作数之后
func (p *RpnPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	var builder strings.Builder

	// 先处理操作数
	builder.WriteString(expr.Right.Accept(p).(string))
	builder.WriteString(" ")

	// 再处理操作符
	builder.WriteString(expr.Operator.Lexeme)

	return builder.String()
}
