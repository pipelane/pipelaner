// Code generated from Pkl module `com.pipelaner.source.components`. DO NOT EDIT.
package components

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Components interface {
}

var _ Components = (*ComponentsImpl)(nil)

type ComponentsImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Components
func LoadFromPath(ctx context.Context, path string) (ret Components, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Components
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Components, error) {
	var ret ComponentsImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
