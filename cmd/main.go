package main

import (
	"fmt"

	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/refx"

	"github.com/hatlonely/tpl/internal/rpcx"
)

var Version string

type Options struct {
	flag.Options

	Type string
	Rpcx rpcx.Options
}

func main() {
	var options Options
	refx.Must(flag.Struct(&options, refx.WithCamelName()))
	refx.Must(flag.Parse(flag.WithJsonVal()))
	if options.Help {
		fmt.Println(flag.Usage())
		return
	}
	if options.Version {
		fmt.Println(Version)
		return
	}

	//tpl, err := rpcx.NewTemplateWithOptions(&rpcx.Options{
	//	Name:    "rpc-tool",
	//	Package: "github.com/hatlonely/demo",
	//	Registry: struct {
	//		Endpoint  string `dft:"docker.io"`
	//		Namespace string
	//	}{
	//		Endpoint:  "docker.io",
	//		Namespace: "hatlonely",
	//	},
	//	GoProxy:   "https://goproxy.cn",
	//	EnableOps: true,
	//	Ops: struct {
	//		EnableHelm  bool
	//		EnableTrace bool
	//		EnableCors  bool
	//		EnableEsLog bool
	//	}{
	//		EnableHelm:  true,
	//		EnableTrace: true,
	//		EnableCors:  false,
	//		EnableEsLog: true,
	//	},
	//})

	var tpl *rpcx.Template
	var err error

	switch options.Type {
	case "rpcx":
		tpl, err = rpcx.NewTemplateWithOptions(&options.Rpcx)
	}

	if err != nil {
		panic(err)
	}

	if err := tpl.Render("tmp"); err != nil {
		panic(err)
	}
}
