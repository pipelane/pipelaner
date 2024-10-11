/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"time"

	"github.com/docker/go-units"
	"github.com/pipelane/pipelaner"
)

const (
	OptBootstrapServers                 = "bootstrap.servers"
	OptGroupID                          = "group.id"
	OptEnableAutoCommit                 = "enable.auto.commit"
	OptCommitIntervalMs                 = "auto.commit.interval.ms"
	OptAutoOffsetReset                  = "auto.offset.reset"
	OptGoEventsChannelEnable            = "go.events.channel.enable"
	OptSessionTimeoutMs                 = "session.timeout.ms"
	OptHeartBeatIntervalMs              = "heartbeat.interval.ms"
	OptBatchNumMessages                 = "batch.num.messages"
	OptSaslMechanism                    = "sasl.mechanism"
	OptSaslUserName                     = "sasl.username"
	OptSaslPassword                     = "sasl.password"
	OptSecurityProtocol                 = "security.protocol"
	OptBatchSize                        = "batch.size"
	OptRetryBackOffMs                   = "retry.backoff.ms"
	OptRetryBackOffMaxMs                = "retry.backoff.max.ms"
	OptRetries                          = "retries"
	OptEnableIdempotence                = "enable.idempotence"
	OptAcks                             = "acks"
	OptMaxInFlightRequestsPerConnection = "max.in.flight.requests.per.connection"
	OptQueueBufferingMaxMessages        = "queue.buffering.max.messages"
	OptDebug                            = "debug"
	OptFetchMaxWaitMs                   = "fetch.wait.max.ms"
	OptFetchMinBytes                    = "fetch.min.bytes"
	OptQueueBufferingMaxMs              = "queue.buffering.max.ms"
	OptLingerMs                         = "linger.ms"
	OptMaxRequestSize                   = "max.request.size"
	SecuritySaslPlainText               = "sasl_plaintext"
	OptMaxPartitionFetchBytes           = "max.partition.fetch.bytes"
)

type Kafka struct {
	SASLEnabled   bool     `pipelane:"sasl_enabled"`
	SASLMechanism string   `pipelane:"sasl_mechanism"`
	SASLUsername  string   `pipelane:"sasl_username"`
	SASLPassword  string   `pipelane:"sasl_password"`
	Brokers       string   `pipelane:"brokers"`
	Version       string   `pipelane:"version"`
	Topics        []string `pipelane:"topics"`
}

type Config struct {
	SchemaRegistry string `pipelane:"schema_registry"`
	pipelaner.Internal
}

type ConsumerConfig struct {
	Kafka
	Config
	AutoCommitEnabled      bool          `pipelane:"auto_commit_enabled"`
	ConsumerGroupID        string        `pipelane:"consumer_group_id"`
	OffsetNewest           bool          `pipelane:"offset_newest"`
	MaxPartitionFetchBytes *string       `pipelane:"max_partition_fetch_bytes"`
	AutoOffsetReset        string        `pipelane:"auto_offset_reset"`
	ReadTopicTimeout       time.Duration `pipelane:"read_topic_timeout"`
}

func (c *ConsumerConfig) GetMaxPartitionFetchBytes() (int, error) {
	if c.MaxPartitionFetchBytes == nil {
		return 52_428_800, nil
	}
	v, err := units.FromHumanSize(*c.MaxPartitionFetchBytes)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

type ProducerConfig struct {
	Kafka
	Config
	MaxRequestSize            *string `pipelane:"max_request_size"`
	LingerMs                  *int    `pipelane:"linger_ms"`
	QueueBufferingMaxMessages *int    `pipelane:"queue_buffering_max_messages"`
	QueueBufferingMaxMs       *int    `pipelane:"queue_buffering_max_ms"`
	BatchSize                 *int    `pipelane:"batch_size"`
	BatchNumMessages          *int    `pipelane:"batch_num_messages"`
}

func (p *ProducerConfig) GetMaxRequestSize() (int64, error) {
	if p.MaxRequestSize == nil {
		return units.FromHumanSize("1MB")
	}
	return units.FromHumanSize(*p.MaxRequestSize)
}
func (p *ProducerConfig) GetLingerMs() int {
	if p.LingerMs == nil {
		return 100
	}
	l := *p.LingerMs
	return l
}
func (p *ProducerConfig) GetQueueBufferingMaxMessages() int {
	if p.QueueBufferingMaxMessages == nil {
		return 1_000_000
	}
	l := *p.QueueBufferingMaxMessages
	return l
}
func (p *ProducerConfig) GetQueueBufferingMaxMs() int {
	if p.QueueBufferingMaxMs == nil {
		return 1_000
	}
	l := *p.QueueBufferingMaxMs
	return l
}

func (p *ProducerConfig) GetBatchSize() int {
	if p.BatchSize == nil {
		return 16_000_000
	}
	l := *p.BatchSize
	return l
}

func (p *ProducerConfig) GetBatchNumMessages() int {
	if p.BatchNumMessages == nil {
		return 100_000
	}
	l := *p.BatchNumMessages
	return l
}
