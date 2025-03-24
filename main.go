package main

import (
	"fmt"
	"os"

	"github.com/aixiasang/goLox/lox"
)

func main() {
	loxInterpreter := lox.NewLox()

	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		// 从文件运行
		err := loxInterpreter.RunFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(74)
		}
	} else {
		// 以交互模式运行
		fmt.Println("🚀 Lox语言解释器 v1.0")
		fmt.Println("👋 输入 'exit()' 或按 Ctrl+D 退出")
		err := loxInterpreter.RunPrompt()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in REPL: %v\n", err)
			os.Exit(74)
		}
	}
}
