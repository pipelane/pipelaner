package kafka

import (
	"errors"
	"time"

	kcfg "github.com/pipelane/pipelaner/source/shared/kafka"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"

	"github.com/pipelane/pipelaner"
)

type Kafka struct {
	cons   *kafka.Consumer
	cfg    kcfg.Config
	logger zerolog.Logger
}

func NewConsumer(cfg kcfg.Config) (*kafka.Consumer, error) {
	cfgMap := kafka.ConfigMap{
		kcfg.OptBootstrapServers:      cfg.Brokers,
		kcfg.OptGroupID:               cfg.ConsumerGroupID,
		kcfg.OptEnableAutoCommit:      cfg.AutoCommitEnabled,
		kcfg.OptCommitIntervalMs:      time.Millisecond * 500,
		kcfg.OptAutoOffsetReset:       cfg.AutoOffsetReset,
		kcfg.OptGoEventsChannelEnable: false,
		kcfg.OptSessionTimeoutMs:      10000,
		kcfg.OptHeartBeatIntervalMs:   1500,
		kcfg.OptBatchNumMessages:      cfg.BatchSize,
	}

	if cfg.SASLEnabled {
		cfgMap[kcfg.OptSaslMechanism] = cfg.SASLMechanism
		cfgMap[kcfg.OptSaslUserName] = cfg.SASLUsername
		cfgMap[kcfg.OptSaslPassword] = cfg.SASLPassword
		cfgMap[kcfg.OptSecurityProtocol] = kcfg.SecuritySaslPlainText
	}

	cons, err := kafka.NewConsumer(&cfgMap)

	if err != nil {
		return nil, err
	}
	return cons, nil
}

func (c *Kafka) Init(ctx *pipelaner.Context) error {
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
	c.logger = pipelaner.NewLogger()
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
