// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package input

import (
	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/common"
	"github.com/pipelane/pipelaner/gen/source/input/autooffsetreset"
	"github.com/pipelane/pipelaner/gen/source/input/strategy"
)

type KafkaConsumer interface {
	Input

	GetKafka() *common.Kafka

	GetAutoCommitEnabled() bool

	GetConsumerGroupID() string

	GetMaxPartitionFetchBytes() *pkl.DataSize

	GetFetchMaxBytes() *pkl.DataSize

	GetAutoOffsetReset() autooffsetreset.AutoOffsetReset

	GetBalancerStrategy() []strategy.Strategy
}

var _ KafkaConsumer = (*KafkaConsumerImpl)(nil)

type KafkaConsumerImpl struct {
	SourceName string `pkl:"sourceName"`

	Kafka *common.Kafka `pkl:"kafka"`

	AutoCommitEnabled bool `pkl:"autoCommitEnabled"`

	ConsumerGroupID string `pkl:"consumerGroupID"`

	MaxPartitionFetchBytes *pkl.DataSize `pkl:"maxPartitionFetchBytes"`

	FetchMaxBytes *pkl.DataSize `pkl:"fetchMaxBytes"`

	AutoOffsetReset autooffsetreset.AutoOffsetReset `pkl:"autoOffsetReset"`

	BalancerStrategy []strategy.Strategy `pkl:"balancerStrategy"`

	Name string `pkl:"name"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *KafkaConsumerImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *KafkaConsumerImpl) GetKafka() *common.Kafka {
	return rcv.Kafka
}

func (rcv *KafkaConsumerImpl) GetAutoCommitEnabled() bool {
	return rcv.AutoCommitEnabled
}

func (rcv *KafkaConsumerImpl) GetConsumerGroupID() string {
	return rcv.ConsumerGroupID
}

func (rcv *KafkaConsumerImpl) GetMaxPartitionFetchBytes() *pkl.DataSize {
	return rcv.MaxPartitionFetchBytes
}

func (rcv *KafkaConsumerImpl) GetFetchMaxBytes() *pkl.DataSize {
	return rcv.FetchMaxBytes
}

func (rcv *KafkaConsumerImpl) GetAutoOffsetReset() autooffsetreset.AutoOffsetReset {
	return rcv.AutoOffsetReset
}

func (rcv *KafkaConsumerImpl) GetBalancerStrategy() []strategy.Strategy {
	return rcv.BalancerStrategy
}

func (rcv *KafkaConsumerImpl) GetName() string {
	return rcv.Name
}

func (rcv *KafkaConsumerImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *KafkaConsumerImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}
