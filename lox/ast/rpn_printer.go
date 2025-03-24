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

// VisitTernaryExpr 访问三元表达式
// 在RPN中，三元运算符可以被视为条件运算，我们用"?:"表示
func (p *RpnPrinter) VisitTernaryExpr(expr *Ternary) interface{} {
	var builder strings.Builder

	// 条件
	builder.WriteString(expr.Condition.Accept(p).(string))
	builder.WriteString(" ")

	// then分支
	builder.WriteString(expr.ThenBranch.Accept(p).(string))
	builder.WriteString(" ")

	// else分支
	builder.WriteString(expr.ElseBranch.Accept(p).(string))
	builder.WriteString(" ")

	// 最后是三元操作符标记
	builder.WriteString("?:")

	return builder.String()
}

// VisitVariableExpr 访问变量表达式
func (p *RpnPrinter) VisitVariableExpr(expr *Variable) interface{} {
	return expr.Name.Lexeme
}

// VisitAssignExpr 访问赋值表达式
func (p *RpnPrinter) VisitAssignExpr(expr *Assign) interface{} {
	var builder strings.Builder

	builder.WriteString(expr.Value.Accept(p).(string))
	builder.WriteString(" ")
	builder.WriteString(expr.Name.Lexeme)
	builder.WriteString(" =")

	return builder.String()
}

// VisitLogicalExpr 访问逻辑表达式
func (p *RpnPrinter) VisitLogicalExpr(expr *Logical) interface{} {
	var builder strings.Builder

	builder.WriteString(expr.Left.Accept(p).(string))
	builder.WriteString(" ")
	builder.WriteString(expr.Right.Accept(p).(string))
	builder.WriteString(" ")
	builder.WriteString(expr.Operator.Lexeme)

	return builder.String()
}

// VisitCallExpr 访问函数调用表达式
func (p *RpnPrinter) VisitCallExpr(expr *Call) interface{} {
	var builder strings.Builder

	// 先添加所有参数
	for _, arg := range expr.Arguments {
		builder.WriteString(fmt.Sprintf("%v ", arg.Accept(p)))
	}

	// 再添加被调用的表达式
	builder.WriteString(fmt.Sprintf("%v ", expr.Callee.Accept(p)))

	// 添加call操作及参数数量
	builder.WriteString(fmt.Sprintf("call(%d)", len(expr.Arguments)))

	return builder.String()
}
