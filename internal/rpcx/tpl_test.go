package rpcx

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTpl(t *testing.T) {
	Convey("TestTpl", t, func() {
		tpl, err := NewTemplateWithOptions(&Options{
			Name: "rpc-demo",
			Registry: struct {
				Endpoint  string `flag:"usage: docker registry endpoint; default: docker.io"`
				Namespace string `flag:"usage: docker registry namespace"`
			}{
				Endpoint:  "docker.io",
				Namespace: "hatlonely",
			},
			Package:   "github.com/hatlonely/rpcx-demo",
			EnableOps: true,
			GoProxy:   "goproxy.cn",
		})

		So(err, ShouldBeNil)
		os.RemoveAll("tmp")
		So(tpl.Render("tmp"), ShouldBeNil)
	})
}
