/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"fmt"

	"github.com/pipelane/pipelaner/gen/source/common/saslmechanism"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/plugin/kzerolog"
)

type Producer struct {
	*kgo.Client
}

func NewProducer(
	cfg sink.KafkaProducer,
	logger *zerolog.Logger,
) (*Producer, error) {
	mSize := cfg.GetBatchNumMessages()
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.GetKafka().Brokers),
		kgo.WithLogger(kzerolog.New(logger)),
		kgo.ProducerLinger(cfg.GetLingerMs().GoDuration()),
		kgo.MaxBufferedRecords(mSize),
	}

	if cfg.GetKafka().SaslEnabled {
		switch *cfg.GetKafka().SaslMechanism {
		case saslmechanism.SCRAMSHA512, saslmechanism.SCRAMSHA256:
			opts = append(opts,
				kgo.SASL(
					scram.Auth{
						User: *cfg.GetKafka().SaslUsername,
						Pass: *cfg.GetKafka().SaslPassword,
					}.AsSha512Mechanism(),
				),
			)
		default:
			return nil, fmt.Errorf("unknown sasl mechanism: %s", *cfg.GetKafka().SaslMechanism)
		}
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &Producer{Client: cl}, nil
}
