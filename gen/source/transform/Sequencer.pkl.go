// Code generated from Pkl module `com.pipelaner.source.transforms`. DO NOT EDIT.
package transform

type Sequencer interface {
	SequencerImp
}

var _ Sequencer = (*SequencerImpl)(nil)

type SequencerImpl struct {
	SourceName string `pkl:"sourceName"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`

	OutputBufferSize uint `pkl:"outputBufferSize"`
}

func (rcv *SequencerImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *SequencerImpl) GetName() string {
	return rcv.Name
}

func (rcv *SequencerImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *SequencerImpl) GetThreads() uint {
	return rcv.Threads
}

func (rcv *SequencerImpl) GetOutputBufferSize() uint {
	return rcv.OutputBufferSize
}
