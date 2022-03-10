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
				Endpoint  string `dft:"docker.io"`
				Namespace string
			}{
				Endpoint:  "docker.io",
				Namespace: "hatlonely",
			},
			GoProxy: "",
		})

		So(err, ShouldBeNil)
		So(tpl.Render(), ShouldBeNil)
	})
}
