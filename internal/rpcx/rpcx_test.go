package rpcx

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTplMk(t *testing.T) {
	Convey("TestTplMk", t, func() {
		tpl, err := NewTemplateWithOptions(&Options{
			Name: "rpc-tool",
			Registry: struct {
				Endpoint  string `flag:"usage: docker registry endpoint; default: docker.io"`
				Namespace string `flag:"usage: docker registry namespace"`
			}{
				Endpoint:  "docker.io",
				Namespace: "hatlonely",
			},
			Package: "github.com/hatlonely/rpcx-demo",
		})

		So(err, ShouldBeNil)
		So(tpl.Render("."), ShouldBeNil)
	})
}
