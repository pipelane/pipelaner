// Code generated from Pkl module `pipelaner.source.transforms`. DO NOT EDIT.
package transform

type Transform interface {
	GetName() string

	GetSourceName() string

	GetInputs() []string

	GetThreads() int

	GetOutputBufferSize() int
}
