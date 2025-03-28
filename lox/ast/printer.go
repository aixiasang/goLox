package ast

import (
	"fmt"
	"strings"
)

// AstPrinter 实现了表达式访问者接口，用于将AST转换为字符串
type AstPrinter struct{}

// NewAstPrinter 创建一个新的AST打印器
func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

// Print 将表达式转换为字符串
func (p *AstPrinter) Print(expr Expr) string {
	return expr.Accept(p).(string)
}

// VisitBinaryExpr 访问二元表达式
func (p *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

// VisitGroupingExpr 访问分组表达式
func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return p.parenthesize("group", expr.Expression)
}

// VisitLiteralExpr 访问字面量
func (p *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

// VisitUnaryExpr 访问一元表达式
func (p *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

// VisitTernaryExpr 访问三元表达式
func (p *AstPrinter) VisitTernaryExpr(expr *Ternary) interface{} {
	return p.parenthesize("?:", expr.Condition, expr.ThenBranch, expr.ElseBranch)
}

// VisitVariableExpr 访问变量表达式
func (p *AstPrinter) VisitVariableExpr(expr *Variable) interface{} {
	return expr.Name.Lexeme
}

// VisitAssignExpr 访问赋值表达式
func (p *AstPrinter) VisitAssignExpr(expr *Assign) interface{} {
	return p.parenthesize2("=", expr.Name.Lexeme, expr.Value)
}

// VisitLogicalExpr 访问逻辑表达式
func (p *AstPrinter) VisitLogicalExpr(expr *Logical) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

// VisitCallExpr 访问函数调用表达式
func (p *AstPrinter) VisitCallExpr(expr *Call) interface{} {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%v", expr.Callee.Accept(p)))
	builder.WriteString("(")

	for i, arg := range expr.Arguments {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", arg.Accept(p)))
	}

	builder.WriteString(")")

	return builder.String()
}

// parenthesize 将表达式转换为带括号的形式
func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(p).(string))
	}

	builder.WriteString(")")

	return builder.String()
}

// parenthesize2 将表达式转换为带括号的形式，第一个参数可以是字符串
func (p *AstPrinter) parenthesize2(name string, arg interface{}, expr Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)
	builder.WriteString(" ")

	if str, ok := arg.(string); ok {
		builder.WriteString(str)
	} else if e, ok := arg.(Expr); ok {
		builder.WriteString(e.Accept(p).(string))
	}

	builder.WriteString(" ")
	builder.WriteString(expr.Accept(p).(string))
	builder.WriteString(")")

	return builder.String()
}
