/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kcfg "github.com/pipelane/pipelaner/source/shared/kafka"
	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/plugin/kzerolog"
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
	bSize, err := cfg.GetBatchSize()
	if err != nil {
		return nil, err
	}

	cfgMap := kafka.ConfigMap{
		kcfg.OptBootstrapServers:          cfg.Brokers,
		kcfg.OptBatchSize:                 bSize,
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

func NewProducerFrance(
	cfg kcfg.ProducerConfig,
	logger *zerolog.Logger,
) (*kgo.Client, error) {
	lms, err := cfg.GetLingerDurationMs()
	if err != nil {
		return nil, err
	}
	mSize := cfg.GetBatchNumMessages()
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers),
		kgo.WithLogger(kzerolog.New(logger)),
		kgo.ProducerLinger(lms),
		kgo.MaxBufferedRecords(mSize),
	}

	if cfg.SASLEnabled {
		switch cfg.SASLMechanism {
		case "SCRAM-SHA-512":
			opts = append(opts,
				kgo.SASL(
					scram.Sha512(func(ctx context.Context) (scram.Auth, error) {
						return scram.Auth{
							User: cfg.SASLUsername,
							Pass: cfg.SASLPassword,
						}, nil
					},
					),
				),
			)
		case "SCRAM-SHA-256":
			opts = append(opts,
				kgo.SASL(
					scram.Sha256(func(ctx context.Context) (scram.Auth, error) {
						return scram.Auth{
							User: cfg.SASLUsername,
							Pass: cfg.SASLPassword,
						}, nil
					},
					),
				),
			)
		}
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return cl, nil
}
