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

# goLox 解释器

goLox是一个[Lox编程语言](http://craftinginterpreters.com/the-lox-language.html)的Go语言实现，这是一个简单但功能完整的脚本语言，包含变量、控制流、函数、闭包等特性。

## 已实现功能

goLox目前实现了以下功能：

1. **基本数据类型**
   - 数值（浮点数）
   - 字符串
   - 布尔值（true/false）
   - nil

2. **变量和表达式**
   - 变量声明和赋值
   - 算术运算（+, -, *, /）
   - 逻辑运算（and, or, !）
   - 比较运算（==, !=, >, >=, <, <=）
   - 三元运算符（condition ? then : else）

3. **控制流**
   - if/else 条件语句
   - while 循环
   - for 循环
   - break 语句

4. **函数**
   - 函数定义和调用
   - 递归
   - 返回值
   - 闭包
   - 高阶函数（函数作为参数和返回值）

5. **输出**
   - print 语句

## 尚未实现功能

goLox当前版本尚未实现以下功能：

1. **面向对象编程**
   - 类定义
   - 实例创建
   - 方法调用
   - 继承

## 使用方法

### 编译

```bash
go build -o goLox.exe main.go
```

### 运行文件

```bash
./goLox.exe path/to/script.lox
```

### 调试模式

goLox提供了调试模式，用于查看解析和执行过程的详细信息。在调试模式下，解释器会输出词法分析、语法分析和执行过程中的各种信息，包括：

- 当前处理的标记（token）
- 标记匹配情况
- 函数参数解析过程
- 未匹配标记的相关信息

启用调试模式有两种方式：

```bash
# 使用长选项
./goLox.exe --debug path/to/script.lox

# 使用短选项
./goLox.exe -d path/to/script.lox
```

注意：调试标志可以放在命令的任何位置，解释器会自动识别并从参数列表中移除。

### 交互式模式 (REPL)

不提供脚本文件时，goLox会启动交互式解释器：

```bash
./goLox.exe
```

在交互式模式中，您可以逐行输入Lox代码并立即查看执行结果。

## 命令行选项

goLox支持以下命令行选项：

- `--debug` 或 `-d`: 启用调试模式，显示解析和执行的详细信息
- 脚本文件路径: 要执行的Lox脚本文件

用法示例：
```bash
# 运行脚本
./goLox.exe script.lox

# 启用调试模式运行脚本
./goLox.exe --debug script.lox
./goLox.exe -d script.lox

# 启动交互式解释器
./goLox.exe

# 启动调试模式的交互式解释器
./goLox.exe --debug
```

## 语法示例

### 变量和表达式

```
var a = 10;
var b = 20;
var c = a + b;
print c;  // 输出 30
```

### 条件语句

```
if (a > b) {
  print "a大于b";
} else if (a < b) {
  print "a小于b";
} else {
  print "a等于b";
}
```

### 循环

```
// while循环
var count = 0;
while (count < 5) {
  print count;
  count = count + 1;
}

// for循环
for (var i = 0; i < 5; i = i + 1) {
  print i;
}
```

### 函数

```
fun add(a, b) {
  return a + b;
}

print add(5, 3);  // 输出 8
```

### 闭包

```
fun makeCounter() {
  var count = 0;
  fun counter() {
    count = count + 1;
    return count;
  }
  return counter;
}

var counter = makeCounter();
print counter();  // 输出 1
print counter();  // 输出 2
```

### break语句

```
var i = 0;
while (i < 10) {
  print i;
  i = i + 1;
  if (i == 5) {
    break;
  }
}
```

### 递归

```
// 斐波那契数列
fun fibonacci(n) {
  if (n <= 1) return n;
  return fibonacci(n - 1) + fibonacci(n - 2);
}

// 阶乘
fun factorial(n) {
  if (n <= 1) return 1;
  return n * factorial(n - 1);
}
```

### 高阶函数

```
// 函数作为参数
fun applyTwice(func, x) {
  return func(func(x));
}

fun double(x) {
  return x * 2;
}

print applyTwice(double, 3);  // 输出 12 (3→6→12)

// 返回函数
fun makeMultiplier(factor) {
  fun multiply(x) {
    return x * factor;
  }
  return multiply;
}

var triple = makeMultiplier(3);
print triple(4);  // 输出 12
```

## 示例程序

项目中包含了多个示例程序，位于`example`目录下：

- `custom_debug_test.lox`: 测试各种基本功能，包括变量、控制流、函数等
- `function_advanced_test.lox`: 测试函数的高级用法，如闭包和高阶函数
- `solution_test.lox`: 简单的函数调用示例
- `advanced_recursion.lox`: 高级递归示例，包括：
  - 斐波那契数列（普通版和跟踪版）
  - 阶乘计算（普通递归和尾递归优化）
  - 幂函数（线性递归和二分优化）
  - 最大公约数
  - 函数工厂
  - 模拟数组的递归实现

## 限制

由于当前版本的限制，goLox不支持类和继承，因此无法运行使用这些功能的代码。 