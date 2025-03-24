package ast

import (
	"testing"

	"github.com/aixiasang/goLox/lox/token"
)

func TestRpnPrinter(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected string
	}{
		{
			name: "简单二元表达式",
			expr: NewBinary(
				NewLiteral(1),
				token.NewToken(token.PLUS, "+", nil, 1),
				NewLiteral(2),
			),
			expected: "1 2 +",
		},
		{
			name: "嵌套二元表达式",
			expr: NewBinary(
				NewBinary(
					NewLiteral(1),
					token.NewToken(token.PLUS, "+", nil, 1),
					NewLiteral(2),
				),
				token.NewToken(token.STAR, "*", nil, 1),
				NewLiteral(3),
			),
			expected: "1 2 + 3 *",
		},
		{
			name: "带括号的表达式",
			expr: NewBinary(
				NewLiteral(1),
				token.NewToken(token.PLUS, "+", nil, 1),
				NewGrouping(
					NewBinary(
						NewLiteral(2),
						token.NewToken(token.STAR, "*", nil, 1),
						NewLiteral(3),
					),
				),
			),
			expected: "1 2 3 * +",
		},
		{
			name: "一元表达式",
			expr: NewUnary(
				token.NewToken(token.MINUS, "-", nil, 1),
				NewLiteral(123),
			),
			expected: "123 -",
		},
		{
			name: "教程示例: (1 + 2) * (4 - 3)",
			expr: NewBinary(
				NewGrouping(
					NewBinary(
						NewLiteral(1),
						token.NewToken(token.PLUS, "+", nil, 1),
						NewLiteral(2),
					),
				),
				token.NewToken(token.STAR, "*", nil, 1),
				NewGrouping(
					NewBinary(
						NewLiteral(4),
						token.NewToken(token.MINUS, "-", nil, 1),
						NewLiteral(3),
					),
				),
			),
			expected: "1 2 + 4 3 - *",
		},
	}

	printer := NewRpnPrinter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := printer.Print(tt.expr)
			if result != tt.expected {
				t.Errorf("RPN打印结果错误。\n期望: %s\n实际: %s", tt.expected, result)
			}
		})
	}
}
