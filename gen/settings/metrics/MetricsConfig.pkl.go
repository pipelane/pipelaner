// Code generated from Pkl module `com.pipelaner.settings.metrics.MetricsConfig`. DO NOT EDIT.
package metrics

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type MetricsConfig struct {
	Host string `pkl:"host"`

	Port int `pkl:"port"`

	ServiceName *string `pkl:"serviceName"`

	Enable bool `pkl:"enable"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a MetricsConfig
func LoadFromPath(ctx context.Context, path string) (ret *MetricsConfig, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a MetricsConfig
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*MetricsConfig, error) {
	var ret MetricsConfig
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
