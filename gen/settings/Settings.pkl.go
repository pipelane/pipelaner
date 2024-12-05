// Code generated from Pkl module `com.pipelaner.settings.settings`. DO NOT EDIT.
package settings

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/settings/healthcheck"
	"github.com/pipelane/pipelaner/gen/settings/logger"
	"github.com/pipelane/pipelaner/gen/settings/metrics"
)

type Settings struct {
	Logger *logger.Config `pkl:"logger"`

	HealthCheck *healthcheck.Config `pkl:"healthCheck"`

	Metrics *metrics.Config `pkl:"metrics"`

	StartGCAfterMessageProcess bool `pkl:"startGCAfterMessageProcess"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Settings
func LoadFromPath(ctx context.Context, path string) (ret *Settings, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Settings
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Settings, error) {
	var ret Settings
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
