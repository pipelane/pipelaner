/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/pipelane/pipelaner"
	kCfg "github.com/pipelane/pipelaner/source/shared/kafka"
)

type KafkaFrance struct {
	logger *zerolog.Logger
	cfg    kCfg.ProducerConfig
	prod   *kgo.Client
}

func init() {
	pipelaner.RegisterSink("kafka_france", &KafkaFrance{})
}

func (k *KafkaFrance) Init(ctx *pipelaner.Context) error {
	l := ctx.Logger()
	k.logger = &l
	err := ctx.LaneItem().Config().ParseExtended(&k.cfg)
	if err != nil {
		return err
	}

	kafkaLogger := ctx.Logger()
	p, err := NewProducerFrance(k.cfg, &kafkaLogger)
	if err != nil {
		return err
	}

	k.prod = p
	return nil
}

func (k *KafkaFrance) write(ctx context.Context, message []byte) {
	for _, topic := range k.cfg.Topics {
		k.prod.Produce(ctx, &kgo.Record{
			Value: message,
			Topic: topic,
		}, func(record *kgo.Record, err error) {
			if err != nil {
				k.logger.Error().Err(err).Msg("failed to produce message")
				k.write(ctx, message)
			}
		})
	}
}

func (k *KafkaFrance) Sink(ctx *pipelaner.Context, val any) {
	var message []byte

	switch v := val.(type) {
	case []byte:
		message = v
	case string:
		message = []byte(v)
	case chan []byte:
		for vls := range v {
			k.Sink(ctx, vls)
		}
		return
	case chan string:
		for vls := range v {
			k.Sink(ctx, vls)
		}
		return
	case chan any:
		for vls := range v {
			k.Sink(ctx, vls)
		}
		return
	default:
		data, err := json.Marshal(val)
		if err != nil {
			k.logger.Error().Err(err).Msgf("marshall val")
			return
		}
		message = data
	}

	k.write(ctx.Context(), message)
}
