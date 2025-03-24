package resolver

import (
	"testing"

	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/token"
)

// OptimizedTestErrorReporter 为优化解析器测试专用的错误报告器
type OptimizedTestErrorReporter struct {
	errors []string
}

func (er *OptimizedTestErrorReporter) ReportError(line int, message string) {
	er.errors = append(er.errors, message)
}

func (er *OptimizedTestErrorReporter) HasError() bool {
	return len(er.errors) > 0
}

func (er *OptimizedTestErrorReporter) ResetError() {
	er.errors = nil
}

func (er *OptimizedTestErrorReporter) HasRuntimeError() bool {
	return false
}

func (er *OptimizedTestErrorReporter) Error(token *token.Token, line int, message string) {
	er.errors = append(er.errors, message)
}

func NewOptimizedTestErrorReporter() *OptimizedTestErrorReporter {
	return &OptimizedTestErrorReporter{
		errors: make([]string, 0),
	}
}

// 测试变量索引解析
func TestOptimizedVariableIndexing(t *testing.T) {
	// 创建测试错误报告器
	errorReporter := NewOptimizedTestErrorReporter()

	// 创建优化的解析器
	resolver := NewOptimizedResolver(errorReporter)

	// 创建作用域
	resolver.beginScope()

	// 声明并定义变量
	name := token.NewToken(token.IDENTIFIER, "x", nil, 1)
	index := resolver.declare(name)
	resolver.define(name)

	// 验证变量被分配了正确的索引
	if index != 0 {
		t.Errorf("首个变量应该被分配索引0，实际是%d", index)
	}

	// 再声明并定义一个变量
	name2 := token.NewToken(token.IDENTIFIER, "y", nil, 1)
	index2 := resolver.declare(name2)
	resolver.define(name2)

	// 验证第二个变量被分配了正确的索引
	if index2 != 1 {
		t.Errorf("第二个变量应该被分配索引1，实际是%d", index2)
	}

	// 创建变量引用
	varExpr := &ast.Variable{Name: name}

	// 解析变量引用
	resolver.resolveLocal(varExpr, name)

	// 获取位置信息
	locations := resolver.GetLocations()

	// 验证变量位置信息正确
	if location, ok := locations[varExpr]; !ok {
		t.Errorf("未找到变量位置信息")
	} else {
		if location.Depth != 0 || location.Index != 0 {
			t.Errorf("变量位置信息错误，预期(0,0)，实际(%d,%d)", location.Depth, location.Index)
		}
	}

	// 结束作用域
	resolver.endScope()
}

// 测试嵌套作用域中的变量索引
func TestOptimizedNestedScopeVariableIndexing(t *testing.T) {
	// 创建测试错误报告器
	errorReporter := NewOptimizedTestErrorReporter()

	// 创建优化的解析器
	resolver := NewOptimizedResolver(errorReporter)

	// 创建外层作用域
	resolver.beginScope()

	// 声明并定义外层变量
	outerName := token.NewToken(token.IDENTIFIER, "outer", nil, 1)
	outerIndex := resolver.declare(outerName)
	resolver.define(outerName)

	// 创建内层作用域
	resolver.beginScope()

	// 声明并定义内层变量
	innerName := token.NewToken(token.IDENTIFIER, "inner", nil, 2)
	innerIndex := resolver.declare(innerName)
	resolver.define(innerName)

	// 创建对外层变量的引用
	outerVarExpr := &ast.Variable{Name: outerName}

	// 创建对内层变量的引用
	innerVarExpr := &ast.Variable{Name: innerName}

	// 解析变量引用
	resolver.resolveLocal(outerVarExpr, outerName)
	resolver.resolveLocal(innerVarExpr, innerName)

	// 获取位置信息
	locations := resolver.GetLocations()

	// 验证外层变量位置信息
	if location, ok := locations[outerVarExpr]; !ok {
		t.Errorf("未找到外层变量位置信息")
	} else {
		// 深度应该是1，因为我们是在内层作用域引用外层变量
		if location.Depth != 1 || location.Index != outerIndex {
			t.Errorf("外层变量位置信息错误，预期(%d,%d)，实际(%d,%d)",
				1, outerIndex, location.Depth, location.Index)
		}
	}

	// 验证内层变量位置信息
	if location, ok := locations[innerVarExpr]; !ok {
		t.Errorf("未找到内层变量位置信息")
	} else {
		// 深度应该是0，因为我们是在当前作用域引用变量
		if location.Depth != 0 || location.Index != innerIndex {
			t.Errorf("内层变量位置信息错误，预期(%d,%d)，实际(%d,%d)",
				0, innerIndex, location.Depth, location.Index)
		}
	}

	// 结束内层作用域
	resolver.endScope()

	// 结束外层作用域
	resolver.endScope()
}

// 测试自引用检测
func TestOptimizedSelfReferenceDetection(t *testing.T) {
	// 创建测试错误报告器
	errorReporter := NewOptimizedTestErrorReporter()

	// 创建优化的解析器
	resolver := NewOptimizedResolver(errorReporter)

	// 创建作用域
	resolver.beginScope()

	// 声明变量但不定义
	name := token.NewToken(token.IDENTIFIER, "x", nil, 1)
	resolver.declare(name)

	// 创建变量引用表达式
	varExpr := &ast.Variable{Name: name}

	// 尝试解析变量引用（应该报告自引用错误）
	resolver.resolveLocal(varExpr, name)

	// 验证是否报告了自引用错误
	if len(errorReporter.errors) == 0 {
		t.Errorf("未检测到自引用错误")
	} else {
		hasError := false
		for _, err := range errorReporter.errors {
			if containsSubstring(err, "不能在变量初始化中引用自身") {
				hasError = true
				break
			}
		}
		if !hasError {
			t.Errorf("错误消息不符合预期，实际错误: %v", errorReporter.errors)
		}
	}

	// 结束作用域
	resolver.endScope()
}

// 辅助函数，检查字符串是否包含子串
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
