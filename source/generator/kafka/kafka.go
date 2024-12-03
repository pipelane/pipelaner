package kafka

import (
	"context"
	"fmt"

	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/twmb/franz-go/pkg/kgo"
)

func init() {
	source.RegisterInput("kafka-consumer", &Kafka{})
}

type Kafka struct {
	components.Logger
	cons *Consumer
	cfg  input.KafkaConsumer
}

func (c *Kafka) Init(cfg input.Input) error {
	consumerCfg, ok := cfg.(input.KafkaConsumer)
	if !ok {
		return fmt.Errorf("invalid cafka config type: %T", cfg)
	}
	l := c.Log().With().Logger()
	cons, err := NewConsumer(consumerCfg, &l)
	if err != nil {
		return err
	}
	c.cons = cons
	c.cfg = consumerCfg
	return nil
}

func (c *Kafka) Generate(ctx context.Context, input chan<- any) {
	l := c.Log()
	for {
		err := c.cons.Consume(ctx, func(record *kgo.Record) error {
			input <- record.Value
			return nil
		})
		if err != nil {
			l.Error().Err(err).Msg("consume error")
		}
	}
}
