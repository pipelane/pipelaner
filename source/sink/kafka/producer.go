/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kcfg "github.com/pipelane/pipelaner/source/shared/kafka"
)

func NewProducer(cfg kcfg.ProducerConfig) (*kafka.Producer, error) {
	maxSize, err := cfg.GetMaxRequestSize()
	if err != nil {
		return nil, err
	}
	qms, err := cfg.GetQueueBufferingMaxMs()
	if err != nil {
		return nil, err
	}
	lms, err := cfg.GetLingerMs()
	if err != nil {
		return nil, err
	}
	cfgMap := kafka.ConfigMap{
		kcfg.OptBootstrapServers:          cfg.Brokers,
		kcfg.OptBatchSize:                 cfg.GetBatchSize(),
		kcfg.OptBatchNumMessages:          cfg.GetBatchNumMessages(),
		"go.batch.producer":               true,
		kcfg.OptQueueBufferingMaxMessages: cfg.GetQueueBufferingMaxMessages(),
		kcfg.OptQueueBufferingMaxMs:       qms,
		kcfg.OptLingerMs:                  lms,
		kcfg.OptMaxRequestSize:            maxSize,
	}

	if cfg.SASLEnabled {
		cfgMap[kcfg.OptSaslMechanism] = cfg.SASLMechanism
		cfgMap[kcfg.OptSaslUserName] = cfg.SASLUsername
		cfgMap[kcfg.OptSaslPassword] = cfg.SASLPassword
		cfgMap[kcfg.OptSecurityProtocol] = kcfg.SecuritySaslPlainText
	}

	cons, err := kafka.NewProducer(&cfgMap)

	if err != nil {
		return nil, err
	}
	return cons, nil
}
