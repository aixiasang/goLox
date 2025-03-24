package environment

import (
	"testing"

	"github.com/aixiasang/goLox/lox/token"
)

func TestEnvironment(t *testing.T) {
	env := NewEnvironment()

	// 测试定义和获取变量
	env.Define("x", 10.0)
	env.Define("y", "hello")

	xToken := &token.Token{Type: token.IDENTIFIER, Lexeme: "x", Line: 1}
	yToken := &token.Token{Type: token.IDENTIFIER, Lexeme: "y", Line: 1}

	if value := env.Get(xToken); value != 10.0 {
		t.Errorf("Expected x to be 10.0, got %v", value)
	}

	if value := env.Get(yToken); value != "hello" {
		t.Errorf("Expected y to be 'hello', got %v", value)
	}

	// 测试变量赋值
	env.Assign(xToken, 42.0)
	if value := env.Get(xToken); value != 42.0 {
		t.Errorf("Expected x to be 42.0 after assignment, got %v", value)
	}
}

func TestNestedEnvironment(t *testing.T) {
	global := NewEnvironment()
	global.Define("x", 10.0)
	global.Define("y", "global")

	local := NewEnclosedEnvironment(global)
	local.Define("z", true)
	local.Define("y", "local") // 覆盖外部环境的y

	xToken := &token.Token{Type: token.IDENTIFIER, Lexeme: "x", Line: 1}
	yToken := &token.Token{Type: token.IDENTIFIER, Lexeme: "y", Line: 1}
	zToken := &token.Token{Type: token.IDENTIFIER, Lexeme: "z", Line: 1}

	// 测试本地环境和全局环境的变量访问
	if value := local.Get(xToken); value != 10.0 {
		t.Errorf("Expected x to be 10.0 from global scope, got %v", value)
	}

	if value := local.Get(yToken); value != "local" {
		t.Errorf("Expected y to be 'local' in local scope, got %v", value)
	}

	if value := local.Get(zToken); value != true {
		t.Errorf("Expected z to be true in local scope, got %v", value)
	}

	// 修改本地环境中的变量
	local.Assign(yToken, "modified")
	if value := local.Get(yToken); value != "modified" {
		t.Errorf("Expected y to be 'modified' after assignment, got %v", value)
	}

	// 确认全局环境中的y没有被修改
	if value := global.Get(yToken); value != "global" {
		t.Errorf("Expected y to still be 'global' in global environment, got %v", value)
	}

	// 修改全局环境中的变量
	local.Assign(xToken, 99.0)
	if value := global.Get(xToken); value != 99.0 {
		t.Errorf("Expected x to be 99.0 in global environment after assignment, got %v", value)
	}
}
