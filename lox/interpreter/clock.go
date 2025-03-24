package interpreter

import (
	"time"
)

// Clock 是一个内置函数，返回自程序启动以来的秒数
type Clock struct{}

// Call 实现Callable接口，返回当前时间的秒数
func (c *Clock) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return float64(time.Now().Unix())
}

// Arity 返回函数参数数量
func (c *Clock) Arity() int {
	return 0
}

// String 返回函数的字符串表示
func (c *Clock) String() string {
	return "<native fn: clock>"
}
