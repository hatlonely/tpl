package rpcx

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

type Options struct {
	Name     string
	Registry struct {
		Endpoint  string `dft:"docker.io"`
		Namespace string
	}
	GoProxy string
}

func NewTemplateWithOptions(options *Options) (*Template, error) {
	tplMk, err := template.New("").Parse(tplMk)
	if err != nil {
		return nil, errors.Wrap(err, "template.New tplMk failed")
	}

	return &Template{
		options: options,
		TplMk:   tplMk,
	}, nil
}

type Template struct {
	options *Options
	TplMk   *template.Template
}

func (t *Template) Render(prefix string) error {
	if err := render(t.TplMk, t.options, prefix, ".tpl.mk"); err != nil {
		return errors.Wrap(err, "render tplMK failed")
	}

	return nil
}

func render(tpl *template.Template, options *Options, prefix string, out string) error {
	prefix, err := filepath.Abs(prefix)
	if err != nil {
		return errors.Wrap(err, "filepath.Abs failed")
	}
	if err := os.MkdirAll(prefix, 0755); err != nil {
		return errors.Wrap(err, "os.MkdirAll failed")
	}
	fp, err := os.Create(fmt.Sprintf("%v/%v", prefix, out))
	if err != nil {
		return errors.Wrap(err, "os.Open failed")
	}
	if err := tpl.Execute(fp, options); err != nil {
		return errors.Wrap(err, "t.TplMk.Execute failed")
	}
	if err := fp.Close(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
