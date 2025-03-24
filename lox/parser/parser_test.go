package parser

import (
	"testing"

	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/token"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []*token.Token
		expected string
	}{
		{
			name: "简单二元表达式",
			tokens: []*token.Token{
				token.NewToken(token.NUMBER, "1", 1.0, 1),
				token.NewToken(token.PLUS, "+", nil, 1),
				token.NewToken(token.NUMBER, "2", 2.0, 1),
				token.NewToken(token.EOF, "", nil, 1),
			},
			expected: "1 2 +",
		},
		{
			name: "复杂表达式: (1 + 2) * 3",
			tokens: []*token.Token{
				token.NewToken(token.LEFT_PAREN, "(", nil, 1),
				token.NewToken(token.NUMBER, "1", 1.0, 1),
				token.NewToken(token.PLUS, "+", nil, 1),
				token.NewToken(token.NUMBER, "2", 2.0, 1),
				token.NewToken(token.RIGHT_PAREN, ")", nil, 1),
				token.NewToken(token.STAR, "*", nil, 1),
				token.NewToken(token.NUMBER, "3", 3.0, 1),
				token.NewToken(token.EOF, "", nil, 1),
			},
			expected: "1 2 + 3 *",
		},
		{
			name: "一元表达式: -123",
			tokens: []*token.Token{
				token.NewToken(token.MINUS, "-", nil, 1),
				token.NewToken(token.NUMBER, "123", 123.0, 1),
				token.NewToken(token.EOF, "", nil, 1),
			},
			expected: "123 -",
		},
		{
			name: "三元表达式: a > b ? c : d",
			tokens: []*token.Token{
				token.NewToken(token.IDENTIFIER, "a", nil, 1),
				token.NewToken(token.GREATER, ">", nil, 1),
				token.NewToken(token.IDENTIFIER, "b", nil, 1),
				token.NewToken(token.QUESTION, "?", nil, 1),
				token.NewToken(token.IDENTIFIER, "c", nil, 1),
				token.NewToken(token.COLON, ":", nil, 1),
				token.NewToken(token.IDENTIFIER, "d", nil, 1),
				token.NewToken(token.EOF, "", nil, 1),
			},
			expected: "a b > c d ?:",
		},
		{
			name: "嵌套三元表达式: a ? b : c ? d : e",
			tokens: []*token.Token{
				token.NewToken(token.IDENTIFIER, "a", nil, 1),
				token.NewToken(token.QUESTION, "?", nil, 1),
				token.NewToken(token.IDENTIFIER, "b", nil, 1),
				token.NewToken(token.COLON, ":", nil, 1),
				token.NewToken(token.IDENTIFIER, "c", nil, 1),
				token.NewToken(token.QUESTION, "?", nil, 1),
				token.NewToken(token.IDENTIFIER, "d", nil, 1),
				token.NewToken(token.COLON, ":", nil, 1),
				token.NewToken(token.IDENTIFIER, "e", nil, 1),
				token.NewToken(token.EOF, "", nil, 1),
			},
			expected: "a b c d e ?: ?:",
		},
		{
			name: "逗号表达式: a, b, c",
			tokens: []*token.Token{
				token.NewToken(token.IDENTIFIER, "a", nil, 1),
				token.NewToken(token.COMMA, ",", nil, 1),
				token.NewToken(token.IDENTIFIER, "b", nil, 1),
				token.NewToken(token.COMMA, ",", nil, 1),
				token.NewToken(token.IDENTIFIER, "c", nil, 1),
				token.NewToken(token.EOF, "", nil, 1),
			},
			expected: "a b c , ,",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := error.NewErrorReporter()
			parser := NewParser(tt.tokens, errors)
			expr := parser.Parse()

			// 使用RPN打印器验证结果
			printer := ast.NewRpnPrinter()
			result := printer.Print(expr)

			if result != tt.expected {
				t.Errorf("解析结果错误。\n期望: %s\n实际: %s", tt.expected, result)
			}

			if errors.HasError() {
				t.Errorf("解析过程中出现错误")
			}
		})
	}
}

func TestParserErrorHandling(t *testing.T) {
	tests := []struct {
		name      string
		tokens    []*token.Token
		expectErr bool
	}{
		{
			name: "缺少左操作数的二元操作符: + 5",
			tokens: []*token.Token{
				token.NewToken(token.PLUS, "+", nil, 1),
				token.NewToken(token.NUMBER, "5", 5.0, 1),
				token.NewToken(token.EOF, "", nil, 1),
			},
			expectErr: true,
		},
		{
			name: "缺少右括号: (1 + 2",
			tokens: []*token.Token{
				token.NewToken(token.LEFT_PAREN, "(", nil, 1),
				token.NewToken(token.NUMBER, "1", 1.0, 1),
				token.NewToken(token.PLUS, "+", nil, 1),
				token.NewToken(token.NUMBER, "2", 2.0, 1),
				token.NewToken(token.EOF, "", nil, 1),
			},
			expectErr: true,
		},
		{
			name: "三元操作符缺少冒号: a ? b",
			tokens: []*token.Token{
				token.NewToken(token.IDENTIFIER, "a", nil, 1),
				token.NewToken(token.QUESTION, "?", nil, 1),
				token.NewToken(token.IDENTIFIER, "b", nil, 1),
				token.NewToken(token.EOF, "", nil, 1),
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := error.NewErrorReporter()
			parser := NewParser(tt.tokens, errors)
			_ = parser.Parse()

			if tt.expectErr && !errors.HasError() {
				t.Errorf("期望解析出错，但没有报告错误")
			}

			if !tt.expectErr && errors.HasError() {
				t.Errorf("期望解析成功，但出现了错误")
			}
		})
	}
}
