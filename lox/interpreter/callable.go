package interpreter

import (
	"github.com/aixiasang/goLox/lox/ast"
	"github.com/aixiasang/goLox/lox/environment"
)

// Callable 表示可调用对象的接口
type Callable interface {
	// Call 调用函数
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	// Arity 返回函数需要的参数数量
	Arity() int
	// String 返回函数的字符串表示
	String() string
}

// Function Lox语言中的函数对象
type Function struct {
	declaration *ast.Function            // 函数声明
	closure     *environment.Environment // 闭包环境
}

// NewFunction 创建一个新的函数对象
func NewFunction(declaration *ast.Function, closure *environment.Environment) *Function {
	return &Function{
		declaration: declaration,
		closure:     closure,
	}
}

// Call 实现Callable接口，调用函数
func (f *Function) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	// 创建函数本地环境，包含参数
	env := environment.NewEnclosedEnvironment(f.closure)

	// 绑定参数
	for i, param := range f.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	// 默认返回值为nil
	var result interface{} = nil

	// 执行函数体并处理return语句（不使用defer）
	func() {
		defer func() {
			if r := recover(); r != nil {
				if returnValue, ok := r.(ReturnValue); ok {
					// 如果是return语句，保存返回值
					result = returnValue.Value
				} else {
					// 其他异常继续传递
					panic(r)
				}
			}
		}()

		// 执行函数体
		interpreter.executeBlock(f.declaration.Body, env)
	}()

	// 返回函数结果，如果没有显式返回则为nil
	return result
}

// Arity 返回函数参数数量
func (f *Function) Arity() int {
	return len(f.declaration.Params)
}

// String 返回函数的字符串表示
func (f *Function) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}

// ReturnValue 用于实现return语句的控制流
type ReturnValue struct {
	Value interface{}
}

// Error 实现error接口
func (r ReturnValue) Error() string {
	return "return"
}
