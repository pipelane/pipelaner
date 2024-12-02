// Code generated from Pkl module `pipelaner.source.example`. DO NOT EDIT.
package custom

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Example struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Example
func LoadFromPath(ctx context.Context, path string) (ret *Example, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Example
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Example, error) {
	var ret Example
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
