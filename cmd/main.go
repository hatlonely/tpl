package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/refx"
	"gopkg.in/yaml.v3"

	"github.com/hatlonely/tpl/internal/tpl"
)

var Version string

type Options struct {
	Help       bool   `flag:"-h; usage: show help info"`
	Version    bool   `flag:"-v; usage: show version"`
	ConfigPath string `flag:"-c; usage: config path"`

	Prefix string `flag:"-p; usage: template directory prefix; default: ."`

	tpl.Options
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

	if _, err := os.Stat(fmt.Sprintf("%v/.tpl.yaml", options.Prefix)); os.IsNotExist(err) {
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf)
		enc.SetIndent(2)
		refx.Must(enc.Encode(&options.Options))
		refx.Must(os.MkdirAll(options.Prefix, 0755))
		refx.Must(ioutil.WriteFile(fmt.Sprintf("%v/.tpl.yaml", options.Prefix), buf.Bytes(), 0644))
	} else {
		buf, err := ioutil.ReadFile(fmt.Sprintf("%v/.tpl.yaml", options.Prefix))
		refx.Must(err)
		refx.Must(yaml.Unmarshal(buf, &options.Options))
	}

	t, err := tpl.NewTemplateWithOptions(&options.Options)
	refx.Must(err)
	refx.Must(t.Render(options.Prefix))
}
