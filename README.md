# 🚀 GoLox - Lox语言的Go实现

这是[Crafting Interpreters](https://craftinginterpreters.fullstack.org.cn/)中Lox语言的Go实现版本。

## ✨ 功能特性

- 📝 完整的词法分析器，支持所有Lox的标记类型
- 💬 支持单行注释 (`//`) 和嵌套的多行注释 (`/* ... */`)
- 🔍 识别所有Lox的关键字、标识符、字符串和数字字面量
- ⚠️ 详细的错误报告，包含行号

## 📂 项目结构

本项目采用模块化结构设计，便于后续扩展开发：

```
.
├── main.go                 # 程序入口点
├── lox/
│   ├── token/              # 标记定义模块
│   │   └── token.go        # Token和TokenType定义
│   ├── scanner/            # 词法分析模块
│   │   └── scanner.go      # 扫描器实现
│   ├── parser/             # 语法分析模块(预留)
│   │   └── parser.go       # 解析器接口定义(预留)
│   ├── ast/                # 抽象语法树模块(预留)
│   │   └── ast.go          # AST节点定义(预留)
│   ├── interpreter/        # 解释执行模块(预留)
│   │   └── interpreter.go  # 解释器实现(预留)
│   ├── error/              # 错误处理模块(预留)
│   │   └── error.go        # 错误处理实现(预留) 
│   └── lox.go              # 解释器主框架
└── test.lox                # 测试文件
```

## 🔧 模块化设计

### 1. Token模块 (`lox/token/`)
- 定义了`TokenType`枚举和`Token`结构
- 提供Token创建和字符串表示功能
- 负责维护关键字和标记类型映射

### 2. Scanner模块 (`lox/scanner/`)
- 实现词法分析器，将源代码转换为标记序列
- 处理字符串、数字、标识符、关键字等
- 支持单行和多行注释

### 3. Parser模块 (`lox/parser/`) - 预留
- 将在后续实现语法分析器
- 将标记序列转换为抽象语法树

### 4. AST模块 (`lox/ast/`) - 预留
- 将定义抽象语法树节点类型
- 支持表达式、语句和声明

### 5. Interpreter模块 (`lox/interpreter/`) - 预留
- 解释执行抽象语法树
- 管理环境和变量作用域

### 6. Error模块 (`lox/error/`) - 预留
- 统一的错误处理机制
- 支持语法错误和运行时错误

## 🚀 使用方法

### 从源代码构建

```bash
# 克隆仓库
git clone https://github.com/your-username/golox.git
cd golox

# 构建项目
go build -o golox
```

### 运行

```bash
# 交互式REPL模式
./golox

# 执行Lox脚本文件
./golox test.lox
```

## 📝 待实现功能

- 🧩 解析器 - 将标记转换为抽象语法树
- ⚙️ 执行器 - 解释和执行抽象语法树
- 🔄 变量绑定和作用域
- 🔀 控制流结构
- 📦 函数和闭包
- 🏗️ 类和继承

## 🌟 扩展功能

本实现比原书中的Java版本添加了一些扩展功能：

1. 🔄 支持嵌套的多行注释
2. 🛡️ 更健壮的边界检查
3. 📦 模块化的代码结构

## 👥 贡献

欢迎通过Pull Request或Issue贡献代码和提出建议。

## 📚 参考资料

- [Crafting Interpreters](https://craftinginterpreters.fullstack.org.cn/) - Lox语言规范和解释器实现指南 