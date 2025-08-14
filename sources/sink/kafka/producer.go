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
	cfg sink.Kafka,
	logger *zerolog.Logger,
) (*Producer, error) {
	mSize := cfg.GetBatchNumMessages()
	ms := cfg.GetLingerMs()

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.GetCommon().Brokers...),
		kgo.WithLogger(kzerolog.New(logger)),
		kgo.ProducerLinger(ms.GoDuration()),
		kgo.MaxBufferedRecords(mSize),
	}

	if cfg.GetTransactionId() != nil {
		opts = append(opts, kgo.TransactionalID(*cfg.GetTransactionId()))
	}

	if cfg.GetCommon().SaslAuth != nil {
		auth := scram.Auth{
			User: cfg.GetCommon().SaslAuth.SaslUsername,
			Pass: cfg.GetCommon().SaslAuth.SaslPassword,
		}
		var authOpt kgo.Opt
		switch cfg.GetCommon().SaslAuth.SaslMechanism {
		case saslmechanism.SCRAMSHA512:
			authOpt = kgo.SASL(auth.AsSha512Mechanism())
		case saslmechanism.SCRAMSHA256:
			authOpt = kgo.SASL(auth.AsSha256Mechanism())
		default:
			return nil, fmt.Errorf("unknown sasl mechanism: %s", cfg.GetCommon().SaslAuth.SaslMechanism)
		}
		opts = append(opts, authOpt)
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &Producer{Client: cl}, nil
}
