/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"

	"github.com/pipelane/pipelaner"
	kCfg "github.com/pipelane/pipelaner/source/shared/kafka"
)

const timeout = 15 * 1000

type Kafka struct {
	logger zerolog.Logger
	cfg    *kCfg.KafkaConfig
	prod   *kafka.Producer
}

func (k *Kafka) Init(ctx *pipelaner.Context) error {
	k.logger = pipelaner.NewLogger()
	k.cfg = new(kCfg.KafkaConfig)
	err := ctx.LaneItem().Config().ParseExtended(k.cfg)
	if err != nil {
		return err
	}

	p, err := NewProducer(k.cfg)
	if err != nil {
		return err
	}

	k.prod = p

	go func() {
		for e := range k.prod.Events() {
			if ev, ok := e.(*kafka.Message); ok {
				if ev.TopicPartition.Error != nil {
					k.logger.Error().Err(ev.TopicPartition.Error).Msgf("delivered failed")
				}
			}
		}
	}()

	return nil
}

func (k *Kafka) write(message []byte) {
	for _, topic := range k.cfg.KafkaTopics {
		if err := k.prod.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          message,
		}, nil); err != nil {
			k.logger.Error().Err(err).Msgf("kafka produce")
			return
		}
	}

	k.prod.Flush(timeout)
}

func (k *Kafka) Sink(_ *pipelaner.Context, val any) {
	var message []byte

	switch v := val.(type) {
	case []byte:
		message = v
	case string:
		message = []byte(v)
	default:
		data, err := json.Marshal(val)
		if err != nil {
			k.logger.Error().Err(err).Msgf("marshall val")
			return
		}

		message = data
	}

	k.write(message)
}
