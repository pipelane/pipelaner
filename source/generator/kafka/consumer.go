package kafka

import (
	"errors"
	kcfg "github.com/pipelane/pipelaner/source/shared/kafka"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pipelane/pipelaner"
	"github.com/rs/zerolog"
)

type Kafka struct {
	cons   *kafka.Consumer
	cfg    *kcfg.KafkaConfig
	logger zerolog.Logger
}

func NewConsumer(cfg *kcfg.KafkaConfig) (*kafka.Consumer, error) {
	cfgMap := kafka.ConfigMap{
		kcfg.OptBootstrapServers:      cfg.KafkaBrokers,
		kcfg.OptGroupID:               cfg.KafkaConsumerGroupId,
		kcfg.OptEnableAutoCommit:      cfg.KafkaAutoCommitEnabled,
		kcfg.OptCommitIntervalMs:      time.Millisecond * 500,
		kcfg.OptAutoOffsetReset:       cfg.KafkaAutoOffsetReset,
		kcfg.OptGoEventsChannelEnable: false,
		kcfg.OptSessionTimeoutMs:      10000,
		kcfg.OptHeartBeatIntervalMs:   1500,
		kcfg.OptBatchNumMessages:      cfg.KafkaBatchSize,
	}

	if cfg.KafkaSASLEnabled {
		cfgMap[kcfg.OptSaslMechanism] = cfg.KafkaSASLMechanism
		cfgMap[kcfg.OptSaslUserName] = cfg.KafkaSASLUsername
		cfgMap[kcfg.OptSaslPassword] = cfg.KafkaSASLPassword
		cfgMap[kcfg.OptSecurityProtocol] = kcfg.SecuritySaslPlainText
	}

	cons, err := kafka.NewConsumer(&cfgMap)

	if err != nil {
		return nil, err
	}
	return cons, nil
}

func (c *Kafka) Init(ctx *pipelaner.Context) error {
	err := ctx.LaneItem().Config().ParseExtended(c.cfg)
	if err != nil {
		return err
	}
	c.cons, err = NewConsumer(c.cfg)
	if err != nil {
		return err
	}
	c.logger = pipelaner.NewLogger()
	err = c.cons.SubscribeTopics(c.cfg.KafkaTopics, nil)
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
