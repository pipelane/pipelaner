// Code generated from Pkl module `pipelaner.source.Common`. DO NOT EDIT.
package common

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Common struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Common
func LoadFromPath(ctx context.Context, path string) (ret *Common, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Common
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Common, error) {
	var ret Common
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
