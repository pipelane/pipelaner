/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"

	"github.com/pipelane/pipelaner"
	kCfg "github.com/pipelane/pipelaner/source/shared/kafka"
)

const timeout = 15 * 1000

type Kafka struct {
	logger *zerolog.Logger
	cfg    kCfg.Config
	prod   *kafka.Producer
}

func init() {
	pipelaner.RegisterSink("kafka", &Kafka{})
}

func (k *Kafka) Init(ctx *pipelaner.Context) error {
	k.logger = ctx.Logger()
	err := ctx.LaneItem().Config().ParseExtended(&k.cfg)
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
	for _, topic := range k.cfg.Topics {
		if err := k.prod.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          message,
		}, nil); err != nil {
			if err.Error() == "Local: Queue full" {
				fmt.Printf("kafka produce: %v\n", err)
				k.flush(timeout)
				k.write(message)
			}
		}
	}
}
func (k *Kafka) flush(time int) {
	for {
		if k.prod.Flush(time) == 0 {
			break
		}
	}
}

func (k *Kafka) Sink(_ *pipelaner.Context, val any) {
	var message []byte

	switch v := val.(type) {
	case []byte:
		message = v
	case string:
		message = []byte(v)
	case chan []byte:
		for vls := range v {
			k.write(vls)
		}
		return
	case chan string:
		for vls := range v {
			k.write([]byte(vls))
		}
		return
	case chan any:
		for vls := range v {
			data, err := json.Marshal(vls)
			if err != nil {
				k.logger.Error().Err(err).Msgf("marshal chan val")
				return
			}
			k.write(data)
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

	k.write(message)
}
