/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kcfg "github.com/pipelane/pipelaner/source/shared/kafka"
)

func NewProducer(cfg kcfg.Config) (*kafka.Producer, error) {
	cfgMap := kafka.ConfigMap{
		kcfg.OptBootstrapServers:     cfg.Brokers,
		kcfg.OptBatchNumMessages:     cfg.BatchSize,
		"go.batch.producer":          "true",
		"delivery.report.only.error": "true",
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
