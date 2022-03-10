package main

import (
	"github.com/hatlonely/tpl/internal/rpcx"
)

func main() {
	tpl, err := rpcx.NewTemplateWithOptions(&rpcx.Options{
		Name: "rpc-tool",
		Registry: struct {
			Endpoint  string `dft:"docker.io"`
			Namespace string
		}{
			Endpoint:  "docker.io",
			Namespace: "hatlonely",
		},
		GoProxy: "",
	})

	if err != nil {
		panic(err)
	}

	if err := tpl.Render("tmp"); err != nil {
		panic(err)
	}
}
