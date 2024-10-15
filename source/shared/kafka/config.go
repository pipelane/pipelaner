/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/pipelane/pipelaner"
)

const (
	OptBootstrapServers          = "bootstrap.servers"
	OptGroupID                   = "group.id"
	OptEnableAutoCommit          = "enable.auto.commit"
	OptCommitIntervalMs          = "auto.commit.interval.ms"
	OptAutoOffsetReset           = "auto.offset.reset"
	OptGoEventsChannelEnable     = "go.events.channel.enable"
	OptSessionTimeoutMs          = "session.timeout.ms"
	OptHeartBeatIntervalMs       = "heartbeat.interval.ms"
	OptBatchNumMessages          = "batch.num.messages"
	OptSaslMechanism             = "sasl.mechanism"
	OptSaslUserName              = "sasl.username"
	OptSaslPassword              = "sasl.password"
	OptSecurityProtocol          = "security.protocol"
	OptBatchSize                 = "batch.size"
	OptQueueBufferingMaxMessages = "queue.buffering.max.messages"
	OptQueueBufferingMaxMs       = "queue.buffering.max.ms"
	OptLingerMs                  = "linger.ms"
	OptMaxRequestSize            = "message.max.bytes" // https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md
	SecuritySaslPlainText        = "sasl_plaintext"
	OptMaxPartitionFetchBytes    = "max.partition.fetch.bytes"
	OptFetchMaxBytes             = "fetch.max.bytes"
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
	SchemaRegistry     string `pipelane:"schema_registry"`
	pipelaner.Internal `pipelane:",squash"`
}

type ConsumerConfig struct {
	Kafka             `pipelane:",squash"`
	Config            `pipelane:",squash"`
	AutoCommitEnabled bool   `pipelane:"auto_commit_enabled"`
	ConsumerGroupID   string `pipelane:"consumer_group_id"`
	OffsetNewest      bool   `pipelane:"offset_newest"`

	MaxPartitionFetchBytes string        `pipelane:"max_partition_fetch_bytes"`
	AutoOffsetReset        string        `pipelane:"auto_offset_reset"`
	ReadTopicTimeout       time.Duration `pipelane:"read_topic_timeout"`
	FetchMaxBytes          string        `pipelane:"fetch_max_bytes"`
}

// GetMaxPartitionFetchBytes "B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"
// OR "B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB".
func (c *ConsumerConfig) GetMaxPartitionFetchBytes() (int, error) {
	if c.MaxPartitionFetchBytes == "" {
		return 52_428_800, nil
	}
	v, err := units.FromHumanSize(c.MaxPartitionFetchBytes)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// GetFetchMaxBytes "B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB" OR "B",
// "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB".
func (c *ConsumerConfig) GetFetchMaxBytes() (int, error) {
	if c.FetchMaxBytes == "" {
		return 104_857_600, nil
	}
	v, err := units.FromHumanSize(c.MaxPartitionFetchBytes)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

type ProducerConfig struct {
	Kafka                     `pipelane:",squash"`
	Config                    `pipelane:",squash"`
	MaxRequestSize            string `pipelane:"max_request_size"`
	LingerMs                  string `pipelane:"linger_ms"`
	QueueBufferingMaxMessages *int   `pipelane:"queue_buffering_max_messages"`
	QueueBufferingMaxMs       string `pipelane:"queue_buffering_max_ms"`
	BatchSize                 *int   `pipelane:"batch_size"`
	BatchNumMessages          *int   `pipelane:"batch_num_messages"`
}

// GetMaxRequestSize "B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB" OR "B",
// "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB".
func (p *ProducerConfig) GetMaxRequestSize() (int, error) {
	if p.MaxRequestSize == "" {
		v, e := units.FromHumanSize("10MiB")
		if e != nil {
			return 0, e
		}
		return int(v), nil
	}
	str := strings.ReplaceAll(p.MaxRequestSize, " ", "")
	v, e := units.FromHumanSize(str)
	if e != nil {
		return 0, e
	}
	return int(v), nil
}

// GetLingerMs "1 ms and etc".
func (p *ProducerConfig) GetLingerMs() (int, error) {
	if p.LingerMs == "" {
		return 100, nil
	}
	l, err := time.ParseDuration(p.LingerMs)
	if err != nil {
		return 0, err
	}
	return int(l.Milliseconds()), nil
}
func (p *ProducerConfig) GetQueueBufferingMaxMessages() int {
	if p.QueueBufferingMaxMessages == nil {
		return 1_000_000
	}
	l := *p.QueueBufferingMaxMessages
	return l
}

// GetQueueBufferingMaxMs "1 ms and etc".
func (p *ProducerConfig) GetQueueBufferingMaxMs() (int, error) {
	if p.QueueBufferingMaxMs == "" {
		return 1_000, nil
	}
	l, err := time.ParseDuration(p.QueueBufferingMaxMs)
	if err != nil {
		return 0, err
	}
	return int(l.Milliseconds()), nil
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
