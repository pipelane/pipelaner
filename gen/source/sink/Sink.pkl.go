// Code generated from Pkl module `pipelaner.source.sinks`. DO NOT EDIT.
package sink

type Sink interface {
	GetName() string

	GetSourceName() string

	GetInputs() []string

	GetThreads() int
}
