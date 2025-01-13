// Code generated from Pkl module `com.pipelaner.source.transforms`. DO NOT EDIT.
package transform

type Transform interface {
	GetName() string

	GetSourceName() string

	GetInputs() []string

	GetThreads() uint

	GetOutputBufferSize() uint
}
