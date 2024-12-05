// Code generated from Pkl module `com.pipelaner.source.components`. DO NOT EDIT.
package components

import (
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/gen/source/transform"
)

type Pipeline struct {
	Name string `pkl:"name"`

	Inputs []input.Input `pkl:"inputs"`

	Transforms []transform.Transform `pkl:"transforms"`

	Sinks []sink.Sink `pkl:"sinks"`
}
