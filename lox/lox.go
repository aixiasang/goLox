package lox

import (
	"bufio"
	"fmt"
	"os"

	errorp "github.com/aixiasang/goLox/lox/error"
	"github.com/aixiasang/goLox/lox/interpreter"
	"github.com/aixiasang/goLox/lox/parser"
	"github.com/aixiasang/goLox/lox/resolver"
	"github.com/aixiasang/goLox/lox/scanner"
)

// Lox 解释器的主结构
type Lox struct {
	errorReporter errorp.Reporter
	interpreter   *interpreter.Interpreter
	debug         bool // 调试模式标志
}

// NewLox 创建一个新的Lox解释器实例
func NewLox() *Lox {
	errorReporter := errorp.NewErrorReporter()
	interpreter := interpreter.NewInterpreter(errorReporter)

	return &Lox{
		errorReporter: errorReporter,
		interpreter:   interpreter,
		debug:         false, // 默认关闭调试模式
	}
}

// SetDebug 设置调试模式
func (l *Lox) SetDebug(debug bool) {
	l.debug = debug
}

// Run 执行给定的源代码
func (l *Lox) Run(source string) {
	// 重置错误状态
	l.errorReporter.ResetError()

	// 扫描标记
	s := scanner.NewScanner(source, l.errorReporter)
	// 设置scanner的调试模式
	s.SetDebug(l.debug)
	tokens := s.ScanTokens()

	// 解析语句
	p := parser.NewParser(tokens, l.errorReporter)
	// 设置解析器的调试模式
	p.SetDebug(l.debug)
	statements := p.Parse()

	// 如果有语法错误,停止解释
	if l.errorReporter.HasError() {
		return
	}

	// 变量解析
	r := resolver.NewResolver(l.interpreter, l.errorReporter)
	r.Resolve(statements)

	// 如果解析过程中有错误,停止解释
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
