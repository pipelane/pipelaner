// Code generated from Pkl module `com.pipelaner.settings.logger.LoggerConfig`. DO NOT EDIT.
package logger

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/settings/logger/logformat"
	"github.com/pipelane/pipelaner/gen/settings/logger/loglevel"
)

type LoggerConfig struct {
	LogLevel loglevel.LogLevel `pkl:"logLevel"`

	EnableConsole bool `pkl:"enableConsole"`

	EnableFile bool `pkl:"enableFile"`

	FileDirectory *string `pkl:"fileDirectory"`

	FileName *string `pkl:"fileName"`

	FileMaxSize *int `pkl:"fileMaxSize"`

	FileMaxBackups *int `pkl:"fileMaxBackups"`

	FileMaxAge *int `pkl:"fileMaxAge"`

	FileCompress *bool `pkl:"fileCompress"`

	FileLocalFormat *bool `pkl:"fileLocalFormat"`

	LogFormat logformat.LogFormat `pkl:"logFormat"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a LoggerConfig
func LoadFromPath(ctx context.Context, path string) (ret *LoggerConfig, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a LoggerConfig
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*LoggerConfig, error) {
	var ret LoggerConfig
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
