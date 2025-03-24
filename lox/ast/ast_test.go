package ast

import (
	"testing"

	"github.com/aixiasang/goLox/lox/token"
)

func TestAstPrinter(t *testing.T) {
	// 创建表达式: -123 * (45.67)
	expression := NewBinary(
		NewUnary(
			token.NewToken(token.MINUS, "-", nil, 1),
			NewLiteral(123),
		),
		token.NewToken(token.STAR, "*", nil, 1),
		NewGrouping(
			NewLiteral(45.67),
		),
	)

	printer := NewAstPrinter()
	result := printer.Print(expression)
	expected := "(* (- 123) (group 45.67))"

	if result != expected {
		t.Errorf("AST打印结果错误。\n期望: %s\n实际: %s", expected, result)
	}
}
