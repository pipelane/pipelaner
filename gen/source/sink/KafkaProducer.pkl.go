// Code generated from Pkl module `pipelaner.source.sinks`. DO NOT EDIT.
package sink

import (
	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/common"
)

type KafkaProducer interface {
	Sink

	GetKafka() *common.Kafka

	GetMaxRequestSize() *pkl.DataSize

	GetLingerMs() *pkl.Duration

	GetBatchNumMessages() int
}

var _ KafkaProducer = (*KafkaProducerImpl)(nil)

type KafkaProducerImpl struct {
	SourceName string `pkl:"sourceName"`

	Kafka *common.Kafka `pkl:"kafka"`

	MaxRequestSize *pkl.DataSize `pkl:"maxRequestSize"`

	LingerMs *pkl.Duration `pkl:"lingerMs"`

	BatchNumMessages int `pkl:"batchNumMessages"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`
}

func (rcv *KafkaProducerImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *KafkaProducerImpl) GetKafka() *common.Kafka {
	return rcv.Kafka
}

func (rcv *KafkaProducerImpl) GetMaxRequestSize() *pkl.DataSize {
	return rcv.MaxRequestSize
}

func (rcv *KafkaProducerImpl) GetLingerMs() *pkl.Duration {
	return rcv.LingerMs
}

func (rcv *KafkaProducerImpl) GetBatchNumMessages() int {
	return rcv.BatchNumMessages
}

func (rcv *KafkaProducerImpl) GetName() string {
	return rcv.Name
}

func (rcv *KafkaProducerImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *KafkaProducerImpl) GetThreads() int {
	return rcv.Threads
}
