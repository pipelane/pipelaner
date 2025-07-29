// Code generated from Pkl module `com.pipelaner`. DO NOT EDIT.
package pipelaner

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/components"
	"github.com/pipelane/pipelaner/gen/settings"
)

type Pipelaner struct {
	Pipelines []components.Pipeline `pkl:"pipelines"`

	Settings settings.Settings `pkl:"settings"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Pipelaner
func LoadFromPath(ctx context.Context, path string) (ret Pipelaner, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return ret, err
	}
	defer func() {
		cerr := evaluator.Close()
		if err == nil {
			err = cerr
		}
	}()
	ret, err = Load(ctx, evaluator, pkl.FileSource(path))
	return ret, err
}

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Pipelaner
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Pipelaner, error) {
	var ret Pipelaner
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
