// Code generated from Pkl module `com.pipelaner.source.inputs`. DO NOT EDIT.
package input

import (
	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/common"
	"github.com/pipelane/pipelaner/gen/source/input/autooffsetreset"
	"github.com/pipelane/pipelaner/gen/source/input/strategy"
)

type Kafka interface {
	Input

	GetCommon() *common.Kafka

	GetAutoCommitEnabled() bool

	GetConsumerGroupID() string

	GetAutoOffsetReset() autooffsetreset.AutoOffsetReset

	GetBalancerStrategy() []strategy.Strategy

	GetMaxPartitionFetchBytes() *pkl.DataSize

	GetFetchMaxBytes() *pkl.DataSize
}

var _ Kafka = (*KafkaImpl)(nil)

type KafkaImpl struct {
	SourceName string `pkl:"sourceName"`

	Common *common.Kafka `pkl:"common"`

	AutoCommitEnabled bool `pkl:"autoCommitEnabled"`

	ConsumerGroupID string `pkl:"consumerGroupID"`

	AutoOffsetReset autooffsetreset.AutoOffsetReset `pkl:"autoOffsetReset"`

	BalancerStrategy []strategy.Strategy `pkl:"balancerStrategy"`

	MaxPartitionFetchBytes *pkl.DataSize `pkl:"maxPartitionFetchBytes"`

	FetchMaxBytes *pkl.DataSize `pkl:"fetchMaxBytes"`

	Name string `pkl:"name"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *KafkaImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *KafkaImpl) GetCommon() *common.Kafka {
	return rcv.Common
}

func (rcv *KafkaImpl) GetAutoCommitEnabled() bool {
	return rcv.AutoCommitEnabled
}

func (rcv *KafkaImpl) GetConsumerGroupID() string {
	return rcv.ConsumerGroupID
}

func (rcv *KafkaImpl) GetAutoOffsetReset() autooffsetreset.AutoOffsetReset {
	return rcv.AutoOffsetReset
}

func (rcv *KafkaImpl) GetBalancerStrategy() []strategy.Strategy {
	return rcv.BalancerStrategy
}

func (rcv *KafkaImpl) GetMaxPartitionFetchBytes() *pkl.DataSize {
	return rcv.MaxPartitionFetchBytes
}

func (rcv *KafkaImpl) GetFetchMaxBytes() *pkl.DataSize {
	return rcv.FetchMaxBytes
}

func (rcv *KafkaImpl) GetName() string {
	return rcv.Name
}

func (rcv *KafkaImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *KafkaImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}
