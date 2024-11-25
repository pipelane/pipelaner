// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package input

type KafkaConsumer interface {
	Input

	GetKafka() *KafkaConfig

	GetAutoCommitEnabled() bool

	GetConsumerGroupID() string

	GetMaxPartitionFetchBytes() string

	GetFetchMaxBytes() string

	GetBalancerStrategy() []string
}

var _ KafkaConsumer = (*KafkaConsumerImpl)(nil)

type KafkaConsumerImpl struct {
	SourceName string `pkl:"sourceName"`

	Kafka *KafkaConfig `pkl:"kafka"`

	AutoCommitEnabled bool `pkl:"autoCommitEnabled"`

	ConsumerGroupID string `pkl:"consumerGroupID"`

	MaxPartitionFetchBytes string `pkl:"maxPartitionFetchBytes"`

	FetchMaxBytes string `pkl:"fetchMaxBytes"`

	BalancerStrategy []string `pkl:"balancerStrategy"`

	Name string `pkl:"name"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *KafkaConsumerImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *KafkaConsumerImpl) GetKafka() *KafkaConfig {
	return rcv.Kafka
}

func (rcv *KafkaConsumerImpl) GetAutoCommitEnabled() bool {
	return rcv.AutoCommitEnabled
}

func (rcv *KafkaConsumerImpl) GetConsumerGroupID() string {
	return rcv.ConsumerGroupID
}

func (rcv *KafkaConsumerImpl) GetMaxPartitionFetchBytes() string {
	return rcv.MaxPartitionFetchBytes
}

func (rcv *KafkaConsumerImpl) GetFetchMaxBytes() string {
	return rcv.FetchMaxBytes
}

func (rcv *KafkaConsumerImpl) GetBalancerStrategy() []string {
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
