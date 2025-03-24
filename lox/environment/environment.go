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
