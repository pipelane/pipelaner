// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package sink

import (
	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/common"
)

type Kafka interface {
	Sink

	GetCommon() *common.Kafka

	GetMaxRequestSize() *pkl.DataSize

	GetLingerMs() *pkl.Duration

	GetBatchNumMessages() int
}

var _ Kafka = (*KafkaImpl)(nil)

type KafkaImpl struct {
	SourceName string `pkl:"sourceName"`

	Common *common.Kafka `pkl:"common"`

	MaxRequestSize *pkl.DataSize `pkl:"maxRequestSize"`

	LingerMs *pkl.Duration `pkl:"lingerMs"`

	BatchNumMessages int `pkl:"batchNumMessages"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`
}

func (rcv *KafkaImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *KafkaImpl) GetCommon() *common.Kafka {
	return rcv.Common
}

func (rcv *KafkaImpl) GetMaxRequestSize() *pkl.DataSize {
	return rcv.MaxRequestSize
}

func (rcv *KafkaImpl) GetLingerMs() *pkl.Duration {
	return rcv.LingerMs
}

func (rcv *KafkaImpl) GetBatchNumMessages() int {
	return rcv.BatchNumMessages
}

func (rcv *KafkaImpl) GetName() string {
	return rcv.Name
}

func (rcv *KafkaImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *KafkaImpl) GetThreads() uint {
	return rcv.Threads
}
