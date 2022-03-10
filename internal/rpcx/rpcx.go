package rpcx

import (
	"os"
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

func (t *Template) Render() error {
	if err := render(t.TplMk, t.options, ".tpl.mk"); err != nil {
		return errors.Wrap(err, "render tplMK failed")
	}

	return nil
}

func render(tpl *template.Template, options *Options, out string) error {
	fp, err := os.Create(out)
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
