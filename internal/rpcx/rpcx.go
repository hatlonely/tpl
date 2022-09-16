package rpcx

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hatlonely/go-kit/strx"
	"github.com/pkg/errors"
)

//go:embed tpl/README.md
var readmeMd string

//go:embed tpl/.rpcx.mk
var rpcxMk string

//go:embed tpl/Dockerfile
var dockerfile string

//go:embed tpl/.gitignore
var gitignore string

//go:embed tpl/api/api.proto
var apiProto string

//go:embed tpl/internal/service/service.go.tpl
var internalServiceService string

//go:embed tpl/cmd/main.go.tpl
var cmdMain string

//go:embed tpl/config/app.json
var configAppJson string

//go:embed tpl/config/base.json
var configBaseJson string

//go:embed tpl/.ops.yaml
var opsYaml string

//go:embed tpl/Makefile
var makefile string

//go:embed tpl/ops/helm/values-adapter.yaml
var opsHelmValuesAdapterYaml string

//go:embed tpl/.rpcx.ops.yaml
var rpcxOpsYaml string

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
	Editable map[string]bool `flag:"usage: recreate if exist"`
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

	if options.Editable == nil {
		options.Editable = map[string]bool{}
	}
	for key, val := range map[string]bool{
		".rpcx.mk":                     true,
		"Dockerfile":                   true,
		"cmd/main.go":                  true,
		"config/base.json":             true,
		"Makefile":                     false,
		".gitignore":                   false,
		"internal/service/service.go":  false,
		"README.md":                    false,
		"config/app.json":              false,
		".ops.yaml":                    false,
		"ops/helm/values-adapter.yaml": false,
		fmt.Sprintf("api/%v.proto", options.Name): false,
	} {
		if _, ok := options.Editable[key]; ok {
			continue
		}
		options.Editable[key] = val
	}

	fmt.Println(options.Editable)

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
			{Tpl: configBaseJson, Out: "config/base.json"},
			{Tpl: configAppJson, Out: "config/app.json"},
			{Tpl: opsYaml, Out: ".ops.yaml", Disable: !options.EnableOps},
			{Tpl: rpcxOpsYaml, Out: ".rpcx.ops.yaml", Disable: !options.EnableOps},
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
	if _, err := os.Stat(out); !os.IsNotExist(err) && !options.Editable[out] {
		fmt.Printf("skip %v\n", out)
		return nil
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

	tpl, err := template.New("").Parse(ts)
	if err != nil {
		return errors.Wrap(err, "template.New failed")
	}
	if err := tpl.Execute(fp, options); err != nil {
		return errors.Wrap(err, "tpl.Execute failed")
	}
	if err := fp.Close(); err != nil {
		return errors.Wrap(err, "close failed")
	}

	fmt.Printf("render %v success\n", out)
	return nil
}
