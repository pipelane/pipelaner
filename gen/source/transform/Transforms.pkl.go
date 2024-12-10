// Code generated from Pkl module `com.pipelaner.source.transforms`. DO NOT EDIT.
package transform

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Transforms struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Transforms
func LoadFromPath(ctx context.Context, path string) (ret *Transforms, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return nil, err
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Transforms
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Transforms, error) {
	var ret Transforms
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
