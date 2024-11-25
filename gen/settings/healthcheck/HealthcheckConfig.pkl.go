// Code generated from Pkl module `com.pipelaner.settings.healthcheck.HealthcheckConfig`. DO NOT EDIT.
package healthcheck

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type HealthcheckConfig struct {
	Host string `pkl:"host"`

	Port int `pkl:"port"`

	Enable bool `pkl:"enable"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a HealthcheckConfig
func LoadFromPath(ctx context.Context, path string) (ret *HealthcheckConfig, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a HealthcheckConfig
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*HealthcheckConfig, error) {
	var ret HealthcheckConfig
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
