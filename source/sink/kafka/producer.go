/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kcfg "github.com/pipelane/pipelaner/source/shared/kafka"
)

func NewProducer(cfg kcfg.KafkaConfig) (*kafka.Producer, error) {
	cfgMap := kafka.ConfigMap{
		kcfg.OptBootstrapServers: cfg.KafkaBrokers,
	}

	if cfg.KafkaSASLEnabled {
		cfgMap[kcfg.OptSaslMechanism] = cfg.KafkaSASLMechanism
		cfgMap[kcfg.OptSaslUserName] = cfg.KafkaSASLUsername
		cfgMap[kcfg.OptSaslPassword] = cfg.KafkaSASLPassword
		cfgMap[kcfg.OptSecurityProtocol] = kcfg.SecuritySaslPlainText
	}

	cons, err := kafka.NewProducer(&cfgMap)

	if err != nil {
		return nil, err
	}
	return cons, nil
}
