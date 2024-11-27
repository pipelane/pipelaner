// Code generated from Pkl module `pipelaner.source.Components`. DO NOT EDIT.
package components

import (
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/gen/source/transform"
)

type Pipeline struct {
	Name string `pkl:"name"`

	Inputs []input.Input `pkl:"inputs"`

	Maps []transform.Transform `pkl:"maps"`

	Sinks []sink.Sink `pkl:"sinks"`
}