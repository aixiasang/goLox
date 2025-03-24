package main

import (
	"fmt"
	"os"

	"github.com/aixiasang/goLox/lox"
)

func main() {
	loxInstance := lox.NewLox()
	args := os.Args[1:]

	// 处理命令行参数
	var scriptPath string
	var debug bool

	// 检查是否有--debug或-d标志
	for i := 0; i < len(args); i++ {
		if args[i] == "--debug" || args[i] == "-d" {
			debug = true
			// 从参数列表中移除debug标志
			if i < len(args)-1 {
				args = append(args[:i], args[i+1:]...)
			} else {
				args = args[:i]
			}
			i-- // 调整索引，因为我们移除了一个元素
		}
	}

	// 设置调试模式
	loxInstance.SetDebug(debug)

	// 检查参数执行文件，否则启动REPL
	if len(args) > 1 {
		fmt.Println("用法: golox [脚本] [--debug/-d]")
		os.Exit(64)
	} else if len(args) == 1 {
		scriptPath = args[0]
		loxInstance.RunFile(scriptPath)
	} else {
		loxInstance.RunPrompt()
	}
}
