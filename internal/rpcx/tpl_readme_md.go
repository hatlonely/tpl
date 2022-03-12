package rpcx

var readmeMd = `
# {{ .Name }}

## 初始化项目

1. 生成 proto 代码

` + "```" + `
make codegen
` + "```" + `

2. 初始化

` + "```" + `
go mod init {{ .Package }}
go mod tidy
go mod vendor
` + "```" + `

## make 命令

1. 代码生成 ` + "`" + `make codegen` + "`" + `
2. 编译 ` + "`" + `make build` + "`" + `
3. 制作镜像 ` + "`" + `make image` + "`" + `
4. 清理 ` + "`" + `make clean` + "`" + `
`
