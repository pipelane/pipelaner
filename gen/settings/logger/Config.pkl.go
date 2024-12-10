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

	EnableFile bool `pkl:"enableFile"`

	FileDirectory *string `pkl:"fileDirectory"`

	FileName *string `pkl:"fileName"`

	FileMaxSize *pkl.DataSize `pkl:"fileMaxSize"`

	FileMaxBackups *int `pkl:"fileMaxBackups"`

	FileMaxAge *int `pkl:"fileMaxAge"`

	FileCompress *bool `pkl:"fileCompress"`

	FileLocalFormat *bool `pkl:"fileLocalFormat"`

	LogFormat logformat.LogFormat `pkl:"logFormat"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Config
func LoadFromPath(ctx context.Context, path string) (ret *Config, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Config
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Config, error) {
	var ret Config
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
