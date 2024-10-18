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
	logger       *zerolog.Logger
	cfg          kCfg.ProducerConfig
	prod         *kafka.Producer
	deliveryChan chan kafka.Event
}

func init() {
	pipelaner.RegisterSink("kafka", &Kafka{})
}

func (k *Kafka) Init(ctx *pipelaner.Context) error {
	l := ctx.Logger()
	k.logger = &l
	err := ctx.LaneItem().Config().ParseExtended(&k.cfg)
	if err != nil {
		return err
	}

	p, err := NewProducer(k.cfg)
	if err != nil {
		return err
	}

	k.prod = p
	k.deliveryChan = make(chan kafka.Event, k.cfg.GetQueueBufferingMaxMessages())
	go func() {
		for e := range k.deliveryChan {
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
		}, k.deliveryChan); err != nil {
			if err.Error() == "Local: Queue full" {
				k.logger.Error().Err(err).Msg("Requeue")
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

func (k *Kafka) Sink(ctx *pipelaner.Context, val any) {
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

	k.write(message)
}
