/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/internal/logger"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
)

func init() {
	source.RegisterSink("kafka-producer", &Kafka{})
}

type Kafka struct {
	logger *zerolog.Logger
	cfg    sink.KafkaProducer
	prod   *Producer
}

func (k *Kafka) Init(cfg sink.Sink) error {
	kafkaCfg, ok := cfg.(sink.KafkaProducer)
	if !ok {
		return fmt.Errorf("invalid kafka-producer config %T", cfg)
	}
	// todo inject node logger
	kafkaLogger := logger.NewLogger()
	p, err := NewProducer(kafkaCfg, &kafkaLogger)
	if err != nil {
		return fmt.Errorf("init kafka producer: %w", err)
	}
	k.prod = p
	k.cfg = kafkaCfg
	return nil
}

func (k *Kafka) write(ctx context.Context, message []byte) {
	for _, topic := range k.cfg.GetKafka().Topics {
		k.prod.Produce(ctx, &kgo.Record{
			Value: message,
			Topic: topic,
		}, func(record *kgo.Record, err error) {
			if err != nil {
				k.logger.Error().Err(err).Msg("failed to produce message")
				k.write(ctx, record.Value)
			}
		})
	}
}

func (k *Kafka) Sink(val any) {
	var message []byte

	switch v := val.(type) {
	case []byte:
		message = v
	case string:
		message = []byte(v)
	case chan []byte:
		for vls := range v {
			k.Sink(vls)
		}
		return
	case chan string:
		for vls := range v {
			k.Sink(vls)
		}
		return
	case chan any:
		for vls := range v {
			k.Sink(vls)
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

	k.write(context.Background(), message)
}
