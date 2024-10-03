package kafka

import (
	"errors"

	kcfg "github.com/pipelane/pipelaner/source/shared/kafka"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"

	"github.com/pipelane/pipelaner"
)

type Kafka struct {
	cons   *kafka.Consumer
	cfg    kcfg.Config
	logger *zerolog.Logger
}

func init() {
	pipelaner.RegisterGenerator("kafka", &Kafka{})
}

func (c *Kafka) Init(ctx *pipelaner.Context) error {
	c.logger = ctx.Logger()
	err := ctx.LaneItem().Config().ParseExtended(&c.cfg)
	if err != nil {
		return err
	}
	if c.cfg.ReadTopicTimeout == 0 {
		c.cfg.ReadTopicTimeout = -1
	}
	c.cons, err = NewConsumer(c.cfg)
	if err != nil {
		return err
	}

	err = c.cons.SubscribeTopics(c.cfg.Topics, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Kafka) Generate(ctx *pipelaner.Context, input chan<- any) {
	for {
		select {
		case <-ctx.Context().Done():
			return
		default:
			msg, err := c.cons.ReadMessage(c.cfg.ReadTopicTimeout)
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
