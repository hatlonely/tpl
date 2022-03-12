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
	Name      string
	Package   string
	Service   string
	EnvPrefix string
	Registry  struct {
		Endpoint  string `dft:"docker.io"`
		Namespace string
	}
	GoProxy string
}

func NewTemplateWithOptions(options *Options) (*Template, error) {
	if options.Service == "" {
		options.Service = strx.PascalName(options.Name)
	}
	if options.EnvPrefix == "" {
		options.EnvPrefix = strx.SnakeNameAllCaps(options.Name)
	}

	return &Template{
		options: options,
		tpls: []struct {
			tpl string
			out string
		}{
			{tpl: rpcxMk, out: ".rpcx.mk"},
			{tpl: makefile, out: "Makefile"},
			{tpl: dockerfile, out: "Dockerfile"},
			{tpl: gitignore, out: ".gitignore"},
			{tpl: apiProto, out: fmt.Sprintf("api/%v.proto", options.Name)},
			{tpl: internalServiceService, out: "internal/service/service.go"},
			{tpl: cmdMain, out: "cmd/main.go"},
			{tpl: readmeMd, out: "README.md"},
			{tpl: ConfigBaseJson, out: "config/base.json"},
			{tpl: ConfigAppJson, out: "config/app.json"},
			{tpl: opsYaml, out: ".ops.yaml"},
		},
	}, nil
}

type Template struct {
	options *Options
	tpls    []struct {
		tpl string
		out string
	}
}

func (t *Template) Render(prefix string) error {
	for _, info := range t.tpls {
		if err := render(info.tpl, t.options, fmt.Sprintf("%v/%v", prefix, info.out)); err != nil {
			return errors.Wrapf(err, "render %v failed", info.out)
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
