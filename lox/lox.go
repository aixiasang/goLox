package lox

import (
	"bufio"
	"fmt"
	"os"

	errorp "github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/interpreter"
	"github.com/aixiasang/goLox/lox/parser"
	"github.com/aixiasang/goLox/lox/scanner"
)

// Lox 解释器的主结构
type Lox struct {
	errorReporter errorp.Reporter
	interpreter   *interpreter.Interpreter
}

// NewLox 创建一个新的Lox解释器实例
func NewLox() *Lox {
	errorReporter := errorp.NewErrorReporter()
	interpreter := interpreter.NewInterpreter(errorReporter)

	return &Lox{
		errorReporter: errorReporter,
		interpreter:   interpreter,
	}
}

// Run 执行给定的源代码
func (l *Lox) Run(source string) {
	// 重置错误状态
	l.errorReporter.ResetError()

	// 扫描标记
	scanner := scanner.NewScanner(source, l.errorReporter)
	tokens := scanner.ScanTokens()

	// 解析语句
	parser := parser.NewParser(tokens, l.errorReporter)
	statements := parser.Parse()

	// 如果有语法错误,停止解释
	if l.errorReporter.HasError() {
		return
	}

	// 解释执行语句
	l.interpreter.Interpret(statements)
}

// RunFile 从文件中读取并执行源代码
func (l *Lox) RunFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	l.Run(string(bytes))

	// 如果有语法错误，返回错误状态
	if l.errorReporter.HasError() {
		os.Exit(65)
	}

	// 如果有运行时错误,返回运行时错误状态
	if l.errorReporter.HasRuntimeError() {
		os.Exit(70)
	}

	return nil
}

// RunPrompt 提供一个交互式的REPL环境
func (l *Lox) RunPrompt() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			// 如果是EOF,优雅退出
			if err.Error() == "EOF" {
				fmt.Println("再见!")
				return nil
			}
			return err
		}

		// Ctrl+D 或EOF会终止循环
		if line == "" {
			break
		}

		l.Run(line)
		// 在REPL中重置错误状态，以便用户可以继续
		l.errorReporter.ResetError()
	}

	return nil
}

// Error 报告错误
func (l *Lox) Error(line int, message string) {
	l.errorReporter.ReportError(line, message)
}
