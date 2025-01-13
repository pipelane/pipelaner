// Code generated from Pkl module `com.pipelaner.settings.logger.config`. DO NOT EDIT.
package logger

import "github.com/apple/pkl-go/pkl"

type WriteFileParams struct {
	Directory string `pkl:"directory"`

	Name string `pkl:"name"`

	MaxSize *pkl.DataSize `pkl:"maxSize"`

	MaxBackups int `pkl:"maxBackups"`

	MaxAge int `pkl:"maxAge"`

	Compress bool `pkl:"compress"`

	LocalFormat bool `pkl:"localFormat"`
}
