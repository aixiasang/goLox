package environment

import (
	"fmt"

	"github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/token"
)

// Environment 存储变量的环境
type Environment struct {
	enclosing *Environment           // 外围环境
	values    map[string]interface{} // 变量值映射
}

// NewEnvironment 创建一个新的环境
func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    make(map[string]interface{}),
	}
}

// NewEnclosedEnvironment 创建一个包含外围环境的新环境
func NewEnclosedEnvironment(enclosing *Environment) *Environment {
	env := NewEnvironment()
	env.enclosing = enclosing
	return env
}

// Define 定义变量（不关心是否已存在）
func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

// Get 获取变量值
func (e *Environment) Get(name *token.Token) interface{} {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}

	// 查找外围环境
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	// 未找到变量
	panic(error.RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("未定义的变量 '%s'。", name.Lexeme),
	})
}

// Assign 给变量赋值（变量必须已存在）
func (e *Environment) Assign(name *token.Token, value interface{}) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}

	// 在外围环境中查找并赋值
	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}

	// 未找到变量
	panic(error.RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("未定义的变量 '%s'。", name.Lexeme),
	})
}

// Ancestor 获取指定深度的环境
func (e *Environment) Ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

// GetAt 从指定深度的环境中获取变量值
func (e *Environment) GetAt(distance int, name string) interface{} {
	// 获取在指定深度的环境，然后从中获取变量
	environment := e.Ancestor(distance)
	if value, ok := environment.values[name]; ok {
		return value
	}

	// 按照正常逻辑这不应该发生，因为解析器已经确定了变量存在
	return nil
}

// AssignAt 在指定深度的环境中赋值
func (e *Environment) AssignAt(distance int, name *token.Token, value interface{}) {
	// 获取在指定深度的环境，然后进行赋值
	environment := e.Ancestor(distance)
	environment.values[name.Lexeme] = value
}
