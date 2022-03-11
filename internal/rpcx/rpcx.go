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
		return nil, errors.Wrap(err, "template.New .tpl.mk failed")
	}

	makefile, err := template.New("").Parse(makefile)
	if err != nil {
		return nil, errors.Wrap(err, "template.New Makefile failed")
	}

	dockerfile, err := template.New("").Parse(dockerfile)
	if err != nil {
		return nil, errors.Wrap(err, "template.New Dockerfile failed")
	}

	gitignore, err := template.New("").Parse(gitignore)
	if err != nil {
		return nil, errors.Wrap(err, "template.New Gitignore failed")
	}

	apiProto, err := template.New("").Parse(apiProto)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("template.New api/%v.proto failed", options.Name))
	}

	internalServiceService, err := template.New("").Parse(internalServiceService)
	if err != nil {
		return nil, errors.Wrap(err, "template.New internal/service/service.go failed")
	}

	return &Template{
		options:                options,
		TplMk:                  tplMk,
		Makefile:               makefile,
		Dockerfile:             dockerfile,
		Gitignore:              gitignore,
		apiProto:               apiProto,
		internalServiceService: internalServiceService,
	}, nil
}

type Template struct {
	options                *Options
	TplMk                  *template.Template
	Makefile               *template.Template
	Dockerfile             *template.Template
	Gitignore              *template.Template
	apiProto               *template.Template
	internalServiceService *template.Template
}

func (t *Template) Render(prefix string) error {
	if err := render(t.TplMk, t.options, fmt.Sprintf("%v/.tpl.mk", prefix)); err != nil {
		return errors.Wrap(err, "render tplMK failed")
	}

	if err := render(t.Makefile, t.options, fmt.Sprintf("%v/Makefile", prefix)); err != nil {
		return errors.Wrap(err, "render Makefile failed")
	}

	if err := render(t.Dockerfile, t.options, fmt.Sprintf("%v/Dockerfile", prefix)); err != nil {
		return errors.Wrap(err, "render Dockerfile failed")
	}

	if err := render(t.Gitignore, t.options, fmt.Sprintf("%v/.gitignore", prefix)); err != nil {
		return errors.Wrap(err, "render .gitignore failed")
	}

	if err := render(t.apiProto, t.options, fmt.Sprintf("%v/api/%v.proto", prefix, t.options.Name)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("render api/%v.proto failed", t.options.Name))
	}

	if err := render(t.internalServiceService, t.options, fmt.Sprintf("%v/internal/service/service.go", prefix)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("render %v/internal/service/service.go failed", prefix))
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
