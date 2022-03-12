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
		return nil, errors.Wrap(err, "template.New .gitignore failed")
	}

	apiProto, err := template.New("").Parse(apiProto)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("template.New api/%v.proto failed", options.Name))
	}

	internalServiceService, err := template.New("").Parse(internalServiceService)
	if err != nil {
		return nil, errors.Wrap(err, "template.New internal/service/service.go failed")
	}

	cmdMain, err := template.New("").Parse(cmdMain)
	if err != nil {
		return nil, errors.Wrap(err, "template.New cmd/main.go failed")
	}

	readmeMd, err := template.New("").Parse(readmeMd)
	if err != nil {
		return nil, errors.Wrap(err, "template.New README.md failed")
	}

	return &Template{
		options:                options,
		tplMk:                  tplMk,
		makefile:               makefile,
		dockerfile:             dockerfile,
		gitignore:              gitignore,
		apiProto:               apiProto,
		internalServiceService: internalServiceService,
		cmdMain:                cmdMain,
		readmeMd:               readmeMd,
	}, nil
}

type Template struct {
	options                *Options
	tplMk                  *template.Template
	makefile               *template.Template
	dockerfile             *template.Template
	gitignore              *template.Template
	apiProto               *template.Template
	internalServiceService *template.Template
	cmdMain                *template.Template
	readmeMd               *template.Template
}

func (t *Template) Render(prefix string) error {
	if err := render(t.tplMk, t.options, fmt.Sprintf("%v/.tpl.mk", prefix)); err != nil {
		return errors.Wrap(err, "render tplMK failed")
	}

	if err := render(t.makefile, t.options, fmt.Sprintf("%v/Makefile", prefix)); err != nil {
		return errors.Wrap(err, "render makefile failed")
	}

	if err := render(t.dockerfile, t.options, fmt.Sprintf("%v/Dockerfile", prefix)); err != nil {
		return errors.Wrap(err, "render dockerfile failed")
	}

	if err := render(t.gitignore, t.options, fmt.Sprintf("%v/.gitignore", prefix)); err != nil {
		return errors.Wrap(err, "render .gitignore failed")
	}

	if err := render(t.apiProto, t.options, fmt.Sprintf("%v/api/%v.proto", prefix, t.options.Name)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("render api/%v.proto failed", t.options.Name))
	}

	if err := render(t.internalServiceService, t.options, fmt.Sprintf("%v/internal/service/service.go", prefix)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("render %v/internal/service/service.go failed", prefix))
	}

	if err := render(t.cmdMain, t.options, fmt.Sprintf("%v/cmd/main.go", prefix)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("render %v/cmd/main.go failed", prefix))
	}

	if err := render(t.readmeMd, t.options, fmt.Sprintf("%v/README.md", prefix)); err != nil {
		return errors.Wrap(err, "render README.md failed")
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
		return errors.Wrap(err, "tpl.Execute failed")
	}
	if err := fp.Close(); err != nil {
		return errors.Wrap(err, "close failed")
	}

	return nil
}
