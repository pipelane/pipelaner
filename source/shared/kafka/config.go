/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"strings"
	"time"

	"github.com/docker/go-units"
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
}

const (
	ConsumerRoundRobinStrategy        = "round-robin"
	ConsumerCooperativeStickyStrategy = "cooperative-sticky"
	ConsumerRangeStrategy             = "range"
)

type ConsumerConfig struct {
	Kafka             `pipelane:",squash"`
	Config            `pipelane:",squash"`
	AutoCommitEnabled bool   `pipelane:"auto_commit_enabled"`
	ConsumerGroupID   string `pipelane:"consumer_group_id"`

	MaxPartitionFetchBytes string   `pipelane:"max_partition_fetch_bytes"`
	AutoOffsetReset        string   `pipelane:"auto_offset_reset"`
	FetchMaxBytes          string   `pipelane:"fetch_max_bytes"`
	BalancerStrategy       []string `pipelane:"balancer_strategy"`
}

// GetMaxPartitionFetchBytes "B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"
// OR "B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB".
func (c *ConsumerConfig) GetMaxPartitionFetchBytes() (int, error) {
	if c.MaxPartitionFetchBytes == "" {
		c.MaxPartitionFetchBytes = "50MiB"
	}
	v, err := units.FromHumanSize(c.MaxPartitionFetchBytes)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func (c *ConsumerConfig) Get() {

}

// GetFetchMaxBytes "B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB" OR "B",
// "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB".
func (c *ConsumerConfig) GetFetchMaxBytes() (int, error) {
	if c.FetchMaxBytes == "" {
		c.FetchMaxBytes = "10MiB"
	}
	v, err := units.FromHumanSize(c.FetchMaxBytes)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

type ProducerConfig struct {
	Kafka            `pipelane:",squash"`
	Config           `pipelane:",squash"`
	MaxRequestSize   string `pipelane:"max_request_size"`
	LingerMs         string `pipelane:"linger_ms"`
	BatchNumMessages *int   `pipelane:"batch_num_messages"`
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
func (p *ProducerConfig) GetLingerDurationMs() (time.Duration, error) {
	if p.LingerMs == "" {
		return 100, nil
	}
	l, err := time.ParseDuration(p.LingerMs)
	if err != nil {
		return 0, err
	}
	return l, nil
}

func (p *ProducerConfig) GetBatchNumMessages() int {
	if p.BatchNumMessages == nil {
		return 100_000
	}
	l := *p.BatchNumMessages
	return l
}
