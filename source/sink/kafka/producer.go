/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	kcfg "github.com/pipelane/pipelaner/source/shared/kafka"
	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/plugin/kzerolog"
)

func NewProducer(
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
					scram.Auth{
						User: cfg.SASLUsername,
						Pass: cfg.SASLPassword,
					}.AsSha512Mechanism(),
				),
			)
		case "SCRAM-SHA-256":
			opts = append(opts,
				kgo.SASL(
					scram.Auth{
						User: cfg.SASLUsername,
						Pass: cfg.SASLPassword,
					}.AsSha256Mechanism(),
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
