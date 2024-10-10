/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"time"

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
	OptLingerMs                         = "linger.ms"
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
	SecuritySaslPlainText               = "sasl_plaintext"
	QueueBufferingMaxMs                 = "queue.buffering.max.ms"
	LingerMs                            = "linger.ms"
)

type Config struct {
	Brokers                   string        `pipelane:"brokers"`
	Version                   string        `pipelane:"version"`
	OffsetNewest              bool          `pipelane:"offset_newest"`
	SASLEnabled               bool          `pipelane:"sasl_enabled"`
	SASLMechanism             string        `pipelane:"sasl_mechanism"`
	SASLUsername              string        `pipelane:"sasl_username"`
	SASLPassword              string        `pipelane:"sasl_password"`
	AutoCommitEnabled         bool          `pipelane:"auto_commit_enabled"`
	ConsumerGroupID           string        `pipelane:"consumer_group_id"`
	Topics                    []string      `pipelane:"topics"`
	AutoOffsetReset           string        `pipelane:"auto_offset_reset"`
	BatchSize                 int           `pipelane:"batch_size"`
	SchemaRegistry            string        `pipelane:"schema_registry"`
	ReadTopicTimeout          time.Duration `pipelane:"read_topic_timeout"`
	QueueBufferingMaxMessages int           `pipelane:"queue_buffering_max_messages"`
	QueueBufferingMaxMs       int           `pipelane:"queue_buffering_max_ms"`
	LingerMs                  int           `pipelane:"linger_ms"`
	pipelaner.Internal
}
