// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package input

import (
	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/input/autooffsetreset"
	"github.com/pipelane/pipelaner/gen/source/input/saslmechanism"
	"github.com/pipelane/pipelaner/gen/source/input/strategy"
)

type KafkaConsumer interface {
	Input

	GetSaslEnabled() *bool

	GetSaslMechanism() *saslmechanism.SASLMechanism

	GetSaslUsername() *string

	GetSaslPassword() *string

	GetBrokers() string

	GetVersion() *string

	GetTopics() []string

	GetSchemaRegistry() *string

	GetAutoCommitEnabled() *bool

	GetConsumerGroupID() string

	GetMaxPartitionFetchBytes() *pkl.DataSize

	GetFetchMaxBytes() *pkl.DataSize

	GetAutoOffsetReset() autooffsetreset.AutoOffsetReset

	GetBalancerStrategy() *[]strategy.Strategy
}

var _ KafkaConsumer = (*KafkaConsumerImpl)(nil)

type KafkaConsumerImpl struct {
	SourceName string `pkl:"sourceName"`

	SaslEnabled *bool `pkl:"saslEnabled"`

	SaslMechanism *saslmechanism.SASLMechanism `pkl:"saslMechanism"`

	SaslUsername *string `pkl:"saslUsername"`

	SaslPassword *string `pkl:"saslPassword"`

	Brokers string `pkl:"brokers"`

	Version *string `pkl:"version"`

	Topics []string `pkl:"topics"`

	SchemaRegistry *string `pkl:"schemaRegistry"`

	AutoCommitEnabled *bool `pkl:"autoCommitEnabled"`

	ConsumerGroupID string `pkl:"consumerGroupID"`

	MaxPartitionFetchBytes *pkl.DataSize `pkl:"maxPartitionFetchBytes"`

	FetchMaxBytes *pkl.DataSize `pkl:"fetchMaxBytes"`

	AutoOffsetReset autooffsetreset.AutoOffsetReset `pkl:"autoOffsetReset"`

	BalancerStrategy *[]strategy.Strategy `pkl:"balancerStrategy"`

	Name string `pkl:"name"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *KafkaConsumerImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *KafkaConsumerImpl) GetSaslEnabled() *bool {
	return rcv.SaslEnabled
}

func (rcv *KafkaConsumerImpl) GetSaslMechanism() *saslmechanism.SASLMechanism {
	return rcv.SaslMechanism
}

func (rcv *KafkaConsumerImpl) GetSaslUsername() *string {
	return rcv.SaslUsername
}

func (rcv *KafkaConsumerImpl) GetSaslPassword() *string {
	return rcv.SaslPassword
}

func (rcv *KafkaConsumerImpl) GetBrokers() string {
	return rcv.Brokers
}

func (rcv *KafkaConsumerImpl) GetVersion() *string {
	return rcv.Version
}

func (rcv *KafkaConsumerImpl) GetTopics() []string {
	return rcv.Topics
}

func (rcv *KafkaConsumerImpl) GetSchemaRegistry() *string {
	return rcv.SchemaRegistry
}

func (rcv *KafkaConsumerImpl) GetAutoCommitEnabled() *bool {
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

func (rcv *KafkaConsumerImpl) GetBalancerStrategy() *[]strategy.Strategy {
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
