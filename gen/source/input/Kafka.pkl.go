// Code generated from Pkl module `com.pipelaner.source.inputs`. DO NOT EDIT.
package input

import (
	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/common"
	"github.com/pipelane/pipelaner/gen/source/input/autooffsetreset"
	"github.com/pipelane/pipelaner/gen/source/input/isolationlevel"
	"github.com/pipelane/pipelaner/gen/source/input/strategy"
)

type Kafka interface {
	Input

	GetCommon() common.Kafka

	GetCommitSrategy() CommitConfiguration

	GetConsumerGroupID() string

	GetAutoOffsetReset() autooffsetreset.AutoOffsetReset

	GetBalancerStrategy() []strategy.Strategy

	GetIsolationLevel() isolationlevel.IsolationLevel

	GetMaxPartitionFetchBytes() pkl.DataSize

	GetFetchMaxBytes() pkl.DataSize
}

var _ Kafka = KafkaImpl{}

type KafkaImpl struct {
	SourceName string `pkl:"sourceName"`

	Common common.Kafka `pkl:"common"`

	CommitSrategy CommitConfiguration `pkl:"commitSrategy"`

	ConsumerGroupID string `pkl:"consumerGroupID"`

	AutoOffsetReset autooffsetreset.AutoOffsetReset `pkl:"autoOffsetReset"`

	BalancerStrategy []strategy.Strategy `pkl:"balancerStrategy"`

	IsolationLevel isolationlevel.IsolationLevel `pkl:"isolationLevel"`

	MaxPartitionFetchBytes pkl.DataSize `pkl:"maxPartitionFetchBytes"`

	FetchMaxBytes pkl.DataSize `pkl:"fetchMaxBytes"`

	Name string `pkl:"name"`

	Threads uint `pkl:"threads"`

	OutputBufferSize uint `pkl:"outputBufferSize"`
}

func (rcv KafkaImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv KafkaImpl) GetCommon() common.Kafka {
	return rcv.Common
}

func (rcv KafkaImpl) GetCommitSrategy() CommitConfiguration {
	return rcv.CommitSrategy
}

func (rcv KafkaImpl) GetConsumerGroupID() string {
	return rcv.ConsumerGroupID
}

func (rcv KafkaImpl) GetAutoOffsetReset() autooffsetreset.AutoOffsetReset {
	return rcv.AutoOffsetReset
}

func (rcv KafkaImpl) GetBalancerStrategy() []strategy.Strategy {
	return rcv.BalancerStrategy
}

func (rcv KafkaImpl) GetIsolationLevel() isolationlevel.IsolationLevel {
	return rcv.IsolationLevel
}

func (rcv KafkaImpl) GetMaxPartitionFetchBytes() pkl.DataSize {
	return rcv.MaxPartitionFetchBytes
}

func (rcv KafkaImpl) GetFetchMaxBytes() pkl.DataSize {
	return rcv.FetchMaxBytes
}

func (rcv KafkaImpl) GetName() string {
	return rcv.Name
}

func (rcv KafkaImpl) GetThreads() uint {
	return rcv.Threads
}

func (rcv KafkaImpl) GetOutputBufferSize() uint {
	return rcv.OutputBufferSize
}
