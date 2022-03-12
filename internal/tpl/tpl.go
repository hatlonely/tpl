package tpl

import (
	"github.com/pkg/errors"

	"github.com/hatlonely/tpl/internal/rpcx"
)

func NewTemplateWithOptions(options *Options) (Template, error) {
	switch options.Type {
	case "rpcx":
		return rpcx.NewTemplateWithOptions(&options.Rpcx)
	}

	return nil, errors.Errorf("unknown template type [%v]", options.Type)
}

type Options struct {
	Type string
	Rpcx rpcx.Options
}

type Template interface {
	Render(prefix string) error
}
