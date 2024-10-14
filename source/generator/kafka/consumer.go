package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kcfg "github.com/pipelane/pipelaner/source/shared/kafka"
)

func NewConsumer(cfg kcfg.ConsumerConfig) (*kafka.Consumer, error) {
	v, err := cfg.GetMaxPartitionFetchBytes()
	if err != nil {
		return nil, err
	}
	maxByteFetch, err := cfg.GetFetchMaxBytes()
	if err != nil {
		return nil, err
	}
	cfgMap := kafka.ConfigMap{
		kcfg.OptBootstrapServers:       cfg.Brokers,
		kcfg.OptGroupID:                cfg.ConsumerGroupID,
		kcfg.OptEnableAutoCommit:       cfg.AutoCommitEnabled,
		kcfg.OptCommitIntervalMs:       time.Millisecond * 500,
		kcfg.OptAutoOffsetReset:        cfg.AutoOffsetReset,
		kcfg.OptGoEventsChannelEnable:  false,
		kcfg.OptSessionTimeoutMs:       10000,
		kcfg.OptHeartBeatIntervalMs:    1500,
		kcfg.OptMaxPartitionFetchBytes: v,
		kcfg.OptFetchMaxBytes:          maxByteFetch,
	}

	if cfg.SASLEnabled {
		cfgMap[kcfg.OptSaslMechanism] = cfg.SASLMechanism
		cfgMap[kcfg.OptSaslUserName] = cfg.SASLUsername
		cfgMap[kcfg.OptSaslPassword] = cfg.SASLPassword
		cfgMap[kcfg.OptSecurityProtocol] = kcfg.SecuritySaslPlainText
	}

	cons, err := kafka.NewConsumer(&cfgMap)

	if err != nil {
		return nil, err
	}
	return cons, nil
}
