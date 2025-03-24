package error

import (
	"fmt"
	"os"

	"github.com/aixiasang/goLox/lox/token"
)

// 错误标记，用于标识是否发生错误
var HadError = false
var HadRuntimeError = false

// Reporter 错误报告接口
type Reporter interface {
	Error(tok *token.Token, line int, message string)
	ReportError(line int, message string)
	ResetError()
}

// ErrorReporter 错误报告实现
type ErrorReporter struct{}

// NewErrorReporter 创建一个新的错误报告器
func NewErrorReporter() *ErrorReporter {
	return &ErrorReporter{}
}

// HasError 返回是否有错误发生
func (r *ErrorReporter) HasError() bool {
	return HadError
}

// HasRuntimeError 返回是否有运行时错误发生
func (r *ErrorReporter) HasRuntimeError() bool {
	return HadRuntimeError
}

// Error 报告错误
func (r *ErrorReporter) Error(tok *token.Token, line int, message string) {
	if tok != nil {
		if tok.Type == token.EOF {
			r.report(tok.Line, "在文件末尾", message)
		} else {
			r.report(tok.Line, fmt.Sprintf("在 '%s'", tok.Lexeme), message)
		}
	} else if line > 0 {
		r.report(line, "", message)
	} else {
		fmt.Fprintf(os.Stderr, "[错误] %s\n", message)
		HadError = true
	}
}

// ReportError 报告一般性错误（不与特定标记关联）
func (r *ErrorReporter) ReportError(line int, message string) {
	r.report(line, "", message)
}

// ResetError 重置错误标记
func (r *ErrorReporter) ResetError() {
	HadError = false
	HadRuntimeError = false
}

// report 报告错误辅助方法
func (r *ErrorReporter) report(line int, where string, message string) {
	if line > 0 {
		fmt.Fprintf(os.Stderr, "[行 %d] 错误 %s: %s\n", line, where, message)
	} else {
		fmt.Fprintf(os.Stderr, "错误 %s: %s\n", where, message)
	}
	HadError = true
}

// RuntimeError 运行时错误类型
type RuntimeError struct {
	Token   *token.Token
	Message string
}

// Error 实现error接口
func (e RuntimeError) Error() string {
	return e.Message
}

// ParseError 解析错误类型
type ParseError struct {
	Token   *token.Token
	Message string
}

// Error 实现error接口
func (e ParseError) Error() string {
	if e.Token.Type == token.EOF {
		return fmt.Sprintf("[行 %d] 错误 在文件末尾: %s", e.Token.Line, e.Message)
	}
	return fmt.Sprintf("[行 %d] 错误 在 '%s': %s", e.Token.Line, e.Token.Lexeme, e.Message)
}
