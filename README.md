# tpl

一个代码模板生成工具

## 安装

### go install

```shell
go install github.com/hatlonely/tpl@v0.0.1
```

### docker 运行

```shell
docker run -i --tty --rm -v $(pwd):/work docker.io/hatlonely/tpl:0.0.1 tpl -h
```

## 快速入门

```shell
tpl --type rpcx --rpcx.name rpc-demo --rpcx.package github.com/hatlonely/rpcx-demo
```
