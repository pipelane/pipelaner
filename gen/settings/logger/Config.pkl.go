// Code generated from Pkl module `com.pipelaner.settings.logger.config`. DO NOT EDIT.
package logger

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/settings/logger/logformat"
	"github.com/pipelane/pipelaner/gen/settings/logger/loglevel"
)

type Config struct {
	LogLevel loglevel.LogLevel `pkl:"logLevel"`

	EnableConsole bool `pkl:"enableConsole"`

	LogFormat logformat.LogFormat `pkl:"logFormat"`

	FileParams *WriteFileParams `pkl:"fileParams"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Config
func LoadFromPath(ctx context.Context, path string) (ret Config, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Config
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Config, error) {
	var ret Config
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
