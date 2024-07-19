/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pipelane/go-kit/config"
	gokit "github.com/pipelane/go-kit/kafka"
	"github.com/pipelane/pipelaner"
	"github.com/rs/zerolog"
)

const timeout = 15 * 1000

type Kafka struct {
	logger zerolog.Logger
	cfg    *pipelaner.KafkaConfig
	prod   *kafka.Producer
}

func NewKafka(logger zerolog.Logger, cfg *pipelaner.KafkaConfig) *Kafka {
	return &Kafka{
		logger: zerolog.Logger{},
		cfg:    cfg,
	}
}

func (k *Kafka) Init(ctx *pipelaner.Context) error {
	k.logger = pipelaner.NewLogger()

	castCfg := pipelaner.CastConfig[*pipelaner.KafkaConfig, config.Kafka](k.cfg)

	p, err := gokit.NewProducer(castCfg)
	if err != nil {
		return err
	}

	k.prod = p

	go func() {
		for e := range k.prod.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
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

	switch val.(type) {
	case []byte:
		message = val.([]byte)
	case string:
		message = []byte(val.(string))
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
