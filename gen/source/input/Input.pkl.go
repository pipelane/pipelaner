// Code generated from Pkl module `com.pipelaner.source.inputs`. DO NOT EDIT.
package input

type Input interface {
	GetName() string

	GetSourceName() string

	GetThreads() uint

	GetOutputBufferSize() uint
}
