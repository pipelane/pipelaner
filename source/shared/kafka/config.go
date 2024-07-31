/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"github.com/pipelane/pipelaner"
	"time"
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
)

type KafkaConfig struct {
	KafkaBrokers           string        `pipelane:"brokers"`
	KafkaVersion           string        `pipelane:"version"`
	KafkaOffsetNewest      bool          `pipelane:"offset_newest"`
	KafkaSASLEnabled       bool          `pipelane:"sasl_enabled"`
	KafkaSASLMechanism     string        `pipelane:"sasl_mechanism"`
	KafkaSASLUsername      string        `pipelane:"sasl_username"`
	KafkaSASLPassword      string        `pipelane:"sasl_password"`
	KafkaAutoCommitEnabled bool          `pipelane:"auto_commit_enabled"`
	KafkaConsumerGroupId   string        `pipelane:"consumer_group_id"`
	KafkaTopics            []string      `pipelane:"topics"`
	KafkaAutoOffsetReset   string        `pipelane:"auto_offset_reset"`
	KafkaBatchSize         int           `pipelane:"batch_size"`
	KafkaSchemaRegistry    string        `pipelane:"schema_registry"`
	ReadTopicTimeout       time.Duration `pipelane:"read_topic_timeout"`
	pipelaner.Internal
}
