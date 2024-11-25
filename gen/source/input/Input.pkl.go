// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package input

type Input interface {
	GetName() string

	GetSourceName() string

	GetThreads() int

	GetOutputBufferSize() int
}
