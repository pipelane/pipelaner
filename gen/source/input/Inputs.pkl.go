// Code generated from Pkl module `com.pipelaner.source.inputs`. DO NOT EDIT.
package input

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Inputs struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Inputs
func LoadFromPath(ctx context.Context, path string) (ret Inputs, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Inputs
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Inputs, error) {
	var ret Inputs
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
