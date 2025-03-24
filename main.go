package main

import (
	"fmt"
	"os"

	"github.com/aixiasang/goLox/lox"
)

func main() {
	loxInterpreter := lox.NewLox()

	// 检查命令行参数
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("用法: golox [脚本]")
		os.Exit(64)
	} else if len(args) == 1 {
		// 运行脚本文件
		err := loxInterpreter.RunFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "无法读取文件: %v\n", err)
			os.Exit(74)
		}
	} else {
		// 启动交互式REPL
		fmt.Println("GoLox解释器 v0.1.0")
		err := loxInterpreter.RunPrompt()
		if err != nil {
			fmt.Fprintf(os.Stderr, "REPL错误: %v\n", err)
			os.Exit(74)
		}
	}
}
