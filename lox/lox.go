package lox

import (
	"bufio"
	"fmt"
	"os"
)

// Lox 解释器的主结构
type Lox struct {
	hadError bool
}

// NewLox 创建一个新的Lox解释器实例
func NewLox() *Lox {
	return &Lox{
		hadError: false,
	}
}

// Run 执行给定的源代码
func (l *Lox) Run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	// 仅用于调试打印所有的token
	for _, token := range tokens {
		fmt.Println(token)
	}

	// 如果有错误，后续的步骤不会执行（后续会实现解析和执行阶段）
}

// RunFile 从文件中读取并执行源代码
func (l *Lox) RunFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	l.Run(string(bytes))

	// 如果有语法错误，返回错误状态
	if l.hadError {
		os.Exit(65)
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
			return err
		}

		// Ctrl+D 或EOF会终止循环
		if line == "" {
			break
		}

		l.Run(line)
		// 在REPL模式下，错误不会终止程序，重置错误状态
		l.hadError = false
	}

	return nil
}

// Error 报告错误
func (l *Lox) Error(line int, message string) {
	l.Report(line, "", message)
}

// Report 报告详细的错误信息
func (l *Lox) Report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}
