package kafka

import (
	"errors"
	"time"

	gokit "github.com/pipelane/go-kit/kafka"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pipelane/go-kit/config"
	"github.com/pipelane/pipelaner"
	"github.com/rs/zerolog"
)

type Kafka struct {
	cons   *kafka.Consumer
	cfg    *pipelaner.KafkaConfig
	logger zerolog.Logger
}

func NewKafka(
	cfg *pipelaner.KafkaConfig,
	logger zerolog.Logger,
) (*Kafka, error) {
	castCfg := pipelaner.CastConfig[*pipelaner.KafkaConfig, config.Kafka](cfg)

	consumer, err := gokit.NewConsumer(castCfg)
	if err != nil {
		return nil, err
	}
	return &Kafka{
		cons:   consumer,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (c *Kafka) Init(_ *pipelaner.Context) error {
	err := c.cons.SubscribeTopics(c.cfg.KafkaTopics, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Kafka) Generate(ctx *pipelaner.Context, input chan<- any) {
	ticker := time.NewTicker(c.cfg.DelayReadTopic)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Context().Done():
			return
		case <-ticker.C:
			msg, err := c.cons.ReadMessage(-1)
			var kafkaErr *kafka.Error
			if err != nil && errors.As(err, &kafkaErr) && kafkaErr.IsTimeout() {
				c.logger.Warn().Err(err).Msg("kafka consume timeout")
				continue
			}
			if err != nil {
				c.logger.Error().Err(err).Msg("failed kafka consume")
				return
			}
			if msg != nil {
				input <- msg.Value
			}
		}
	}
}
