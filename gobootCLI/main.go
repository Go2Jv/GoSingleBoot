package main

import (
	"fmt"
	"gobootCLI/cli"
	"os"
)

// gobootCLI 是 GoSingleBoot 的项目模板 CLI。
// 运行方式:
//
//	go run ./CLI
//
// 或者编译后运行:
//
//	go build -o goboot .
//	./goboot
//
// 交互式创建一个基于 GoSingleBoot 模板的新项目:
// 询问项目名 -> clone 模板 -> 改写 go.mod 与 import 路径 -> go mod tidy
// -> 清理模板文件 -> 可选 git init。
func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
