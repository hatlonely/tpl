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

	makefile, err := template.New("").Parse(makefile)
	if err != nil {
		return nil, errors.Wrap(err, "template.New makefile failed")
	}

	return &Template{
		options:  options,
		TplMk:    tplMk,
		Makefile: makefile,
	}, nil
}

type Template struct {
	options  *Options
	TplMk    *template.Template
	Makefile *template.Template
}

func (t *Template) Render(prefix string) error {
	if err := render(t.TplMk, t.options, fmt.Sprintf("%v/.tpl.mk", prefix)); err != nil {
		return errors.Wrap(err, "render tplMK failed")
	}

	if err := render(t.TplMk, t.options, fmt.Sprintf("%v/Makefile", prefix)); err != nil {
		return errors.Wrap(err, "render Makefile failed")
	}

	return nil
}

func render(tpl *template.Template, options *Options, out string) error {
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
		return errors.Wrap(err, "t.TplMk.Execute failed")
	}
	if err := fp.Close(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
