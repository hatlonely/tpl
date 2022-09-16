package render

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

type Desc struct {
	Tpl     string
	Out     string
	Disable bool
}

func Render(prefix string, tpls []Desc, options interface{}, editable map[string]bool) error {
	for _, desc := range tpls {
		if desc.Disable {
			continue
		}
		if err := render(desc.Tpl, options, fmt.Sprintf("%v/%v", prefix, desc.Out), editable); err != nil {
			return errors.Wrapf(err, "render %v failed", desc.Out)
		}
	}

	return nil
}

func render(ts string, options interface{}, out string, editable map[string]bool) error {
	if _, err := os.Stat(out); !os.IsNotExist(err) && !editable[out] {
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
