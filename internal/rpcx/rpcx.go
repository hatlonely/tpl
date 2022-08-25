package rpcx

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hatlonely/go-kit/strx"
	"github.com/pkg/errors"
)

type Options struct {
	Name      string `flag:"usage: project name"`
	Package   string `flag:"usage: package name"`
	Service   string `flag:"usage: service name, use pascal Name if not specific"`
	EnvPrefix string `flag:"usage: environment prefix, use all caps snake Name if not specific"`
	ImageType string `flag:"usage: image type, centos|alpine; default: alpine"`
	Registry  struct {
		Endpoint  string `flag:"usage: docker registry endpoint; default: docker.io"`
		Namespace string `flag:"usage: docker registry namespace"`
	}
	GoProxy   string `flag:"usage: set go proxy in Makefile"`
	EnableOps bool   `flag:"usage: generate .ops.yaml"`
	Ops       struct {
		EnableHelm  bool `flag:"usage: enable helm task"`
		EnableTrace bool
		EnableCors  bool
		EnableEsLog bool
	}
}

func NewTemplateWithOptions(options *Options) (*Template, error) {
	if options.Name == "" {
		return nil, errors.New("miss required field [Name]")
	}
	if options.Package == "" {
		return nil, errors.New("miss required field [Package]")
	}
	if options.Service == "" {
		options.Service = strx.PascalName(options.Name)
	}
	if options.EnvPrefix == "" {
		options.EnvPrefix = strx.SnakeNameAllCaps(options.Name)
	}

	return &Template{
		options: options,
		tpls: []TplDesc{
			{Tpl: rpcxMk, Out: ".rpcx.mk"},
			{Tpl: makefile, Out: "Makefile"},
			{Tpl: dockerfile, Out: "Dockerfile"},
			{Tpl: gitignore, Out: ".gitignore"},
			{Tpl: apiProto, Out: fmt.Sprintf("api/%v.proto", options.Name)},
			{Tpl: internalServiceService, Out: "internal/service/service.go"},
			{Tpl: cmdMain, Out: "cmd/main.go"},
			{Tpl: readmeMd, Out: "README.md"},
			{Tpl: ConfigBaseJson, Out: "config/base.json"},
			{Tpl: ConfigAppJson, Out: "config/app.json"},
			{Tpl: opsYaml, Out: ".ops.yaml", Disable: !options.EnableOps},
			{Tpl: rpcxOpsYaml, Out: ".rpcx.ops.yaml"},
			{Tpl: opsHelmValuesAdapterYaml, Out: "ops/helm/values-adapter.yaml", Disable: !options.EnableOps || !options.Ops.EnableHelm},
		},
	}, nil
}

type TplDesc struct {
	Tpl     string
	Out     string
	Disable bool
}

type Template struct {
	options *Options
	tpls    []TplDesc
}

func (t *Template) Render(prefix string) error {
	for _, desc := range t.tpls {
		if desc.Disable {
			continue
		}
		if err := render(desc.Tpl, t.options, fmt.Sprintf("%v/%v", prefix, desc.Out)); err != nil {
			return errors.Wrapf(err, "render %v failed", desc.Out)
		}
	}

	return nil
}

func render(ts string, options *Options, out string) error {
	tpl, err := template.New("").Parse(ts)
	if err != nil {
		return errors.Wrap(err, "template.New failed")
	}

	abs, err := filepath.Abs(out)
	if err != nil {
		return errors.Wrap(err, "filepath.Abs failed")
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		return errors.Wrap(err, "os.MkdirAll failed")
	}
	fp, err := os.Create(abs)
	if err != nil {
		return errors.Wrap(err, "os.Open failed")
	}
	if err := tpl.Execute(fp, options); err != nil {
		return errors.Wrap(err, "tpl.Execute failed")
	}
	if err := fp.Close(); err != nil {
		return errors.Wrap(err, "close failed")
	}

	return nil
}
