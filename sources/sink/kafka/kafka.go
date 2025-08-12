/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/node"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/twmb/franz-go/pkg/kgo"
)

func init() {
	source.RegisterSink("kafka", &Kafka{})
}

type producer interface {
	Produce(
		ctx context.Context,
		r *kgo.Record,
		promise func(*kgo.Record, error),
	)
}

type Kafka struct {
	components.Logger
	cfg  sink.Kafka
	prod producer
}

func (k *Kafka) Init(cfg sink.Sink) error {
	kafkaCfg, ok := cfg.(sink.Kafka)
	if !ok {
		return fmt.Errorf("invalid kafka-producer config %T", cfg)
	}
	kafkaLogger := k.Log().With().Logger()
	p, err := NewProducer(kafkaCfg, &kafkaLogger)
	if err != nil {
		return fmt.Errorf("init kafka producer: %w", err)
	}
	k.prod = p
	k.cfg = kafkaCfg
	return nil
}

func (k *Kafka) write(ctx context.Context, message []byte) {
	for _, topic := range k.cfg.GetCommon().Topics {
		k.prod.Produce(ctx, &kgo.Record{
			Value: message,
			Topic: topic,
		}, func(record *kgo.Record, err error) {
			if err != nil {
				k.Log().Error().Err(err).Msg("failed to produce message")
				k.write(ctx, record.Value)
			}
		})
	}
}

func (k *Kafka) Sink(val any) error {
	var message []byte
	var err error

	switch v := val.(type) {
	case node.AtomicData:
		err = k.Sink(v.Data())
		if err != nil {
			v.Error() <- v
			return err
		}
		v.Success() <- v
		return nil
	case []byte:
		message = v
	case string:
		message = []byte(v)
	case chan node.AtomicData:
		for msg := range v {
			_ = k.Sink(msg) //nolint: errcheck
		}
		return nil
	case chan []byte:
		for msg := range v {
			_ = k.Sink(msg) //nolint: errcheck
		}
		return nil

	case chan string:
		for msg := range v {
			_ = k.Sink(msg) //nolint: errcheck
		}
		return nil

	case chan any:
		for msg := range v {
			_ = k.Sink(msg) //nolint: errcheck
		}
		return nil
	default:
		data, errs := json.Marshal(val)
		if errs != nil {
			k.Log().Error().Err(errs).Msgf("marshal val")
			return errs
		}
		message = data
	}
	k.write(context.Background(), message)
	return nil
}
