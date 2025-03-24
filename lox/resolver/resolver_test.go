package resolver

import (
	"testing"

	"github.com/aixiasang/goLox/lox/token"
)

// 测试错误报告器
type TestErrorReporter struct {
	errors []string
}

func (er *TestErrorReporter) ReportError(line int, message string) {
	er.errors = append(er.errors, message)
}

func (er *TestErrorReporter) HasError() bool {
	return len(er.errors) > 0
}

func (er *TestErrorReporter) ResetError() {
	er.errors = nil
}

func (er *TestErrorReporter) HasRuntimeError() bool {
	return false
}

func (er *TestErrorReporter) Error(token *token.Token, line int, message string) {
	er.errors = append(er.errors, message)
}

func NewTestErrorReporter() *TestErrorReporter {
	return &TestErrorReporter{
		errors: make([]string, 0),
	}
}

// MockResolver 为了测试创建的简化解析器
type MockResolver struct {
	scopes []map[string]bool
	errors *TestErrorReporter
}

func NewMockResolver() *MockResolver {
	return &MockResolver{
		scopes: make([]map[string]bool, 0),
		errors: NewTestErrorReporter(),
	}
}

func (r *MockResolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *MockResolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *MockResolver) declare(name *token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		r.errors.Error(name, 0, "变量 '"+name.Lexeme+"' 已在此作用域中声明。")
	}

	scope[name.Lexeme] = false
}

func (r *MockResolver) define(name *token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true
}

func (r *MockResolver) resolveVariable(name *token.Token) bool {
	if len(r.scopes) == 0 {
		return false
	}

	scope := r.scopes[len(r.scopes)-1]
	if initialized, exists := scope[name.Lexeme]; exists && !initialized {
		r.errors.Error(name, 0, "不能在变量初始化中引用自身。")
		return false
	}
	return true
}

// 测试变量自引用检测
func TestSelfReferenceDetection(t *testing.T) {
	resolver := NewMockResolver()
	resolver.beginScope()

	// 创建变量名
	name := token.NewToken(token.IDENTIFIER, "x", nil, 1)

	// 声明变量但还未初始化完成
	resolver.declare(name)

	// 尝试在初始化前使用变量（模拟 var x = x）
	result := resolver.resolveVariable(name)

	// 验证错误被捕获
	if result {
		t.Errorf("应该检测到变量自引用错误")
	}

	if len(resolver.errors.errors) == 0 {
		t.Errorf("自引用错误未被检测到")
	} else if !contains(resolver.errors.errors, "不能在变量初始化中引用自身") {
		t.Errorf("预期的错误消息未出现，实际错误: %v", resolver.errors.errors)
	}
}

// 测试嵌套作用域
func TestNestedScopes(t *testing.T) {
	resolver := NewMockResolver()

	// 创建外层作用域并声明变量x
	resolver.beginScope()
	outerName := token.NewToken(token.IDENTIFIER, "x", nil, 1)
	resolver.declare(outerName)
	resolver.define(outerName)

	// 创建内层作用域并声明同名变量x
	resolver.beginScope()
	innerName := token.NewToken(token.IDENTIFIER, "x", nil, 2)
	resolver.declare(innerName)
	resolver.define(innerName)

	// 验证没有报错（不同作用域可以有同名变量）
	if len(resolver.errors.errors) > 0 {
		t.Errorf("不应该有错误，但收到: %v", resolver.errors.errors)
	}

	// 验证内层作用域中可以正确引用内层变量
	if !resolver.resolveVariable(innerName) {
		t.Errorf("内层作用域中应该能够引用内层变量")
	}

	// 结束内层作用域
	resolver.endScope()

	// 验证外层作用域中可以正确引用外层变量
	if !resolver.resolveVariable(outerName) {
		t.Errorf("外层作用域中应该能够引用外层变量")
	}
}

// 辅助函数，检查字符串切片是否包含子串
func contains(slice []string, substr string) bool {
	for _, s := range slice {
		if s == substr || contains_substr(s, substr) {
			return true
		}
	}
	return false
}

// 辅助函数，检查字符串是否包含子串
func contains_substr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
