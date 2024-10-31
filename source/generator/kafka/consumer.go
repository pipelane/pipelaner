package kafka

import (
	"context"
	"fmt"
	"time"

	kcfg "github.com/pipelane/pipelaner/source/shared/kafka"
	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/plugin/kzerolog"
)

type Consumer struct {
	cli *kgo.Client
}

func NewConsumer(
	cfg kcfg.ConsumerConfig,
	logger *zerolog.Logger,
) (*Consumer, error) {
	v, err := cfg.GetMaxPartitionFetchBytes()
	if err != nil {
		return nil, err
	}
	maxByteFetch, err := cfg.GetFetchMaxBytes()
	if err != nil {
		return nil, err
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers),
		kgo.WithLogger(kzerolog.New(logger)),
		kgo.ConsumerGroup(cfg.ConsumerGroupID),
		kgo.ConsumeTopics(cfg.Topics...),
		kgo.FetchMaxBytes(int32(maxByteFetch)), //nolint: gosec
		kgo.FetchMaxPartitionBytes(int32(v)),   //nolint: gosec
		kgo.HeartbeatInterval(time.Second),
	}
	var balancers []kgo.GroupBalancer
	for _, s := range cfg.BalancerStrategy {
		switch s {
		case kcfg.ConsumerRangeStrategy:
			balancers = append(balancers, kgo.RangeBalancer())
		case kcfg.ConsumerRoundRobinStrategy:
			balancers = append(balancers, kgo.RoundRobinBalancer())
		case kcfg.ConsumerCooperativeStickyStrategy:
			balancers = append(balancers, kgo.RangeBalancer())
		}
	}
	if len(balancers) == 0 {
		balancers = append(balancers, kgo.RoundRobinBalancer())
	}
	opts = append(opts, kgo.Balancers(balancers...))

	if !cfg.AutoCommitEnabled {
		opts = append(opts, kgo.DisableAutoCommit())
	}

	if cfg.AutoOffsetReset == "earliest" || cfg.AutoOffsetReset == "" {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()))
	} else if cfg.AutoOffsetReset == "latest" {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()))
	}

	if cfg.SASLEnabled {
		switch cfg.SASLMechanism {
		case "SCRAM-SHA-512":
			opts = append(opts,
				kgo.SASL(
					scram.Auth{
						User: cfg.SASLUsername,
						Pass: cfg.SASLPassword,
					}.AsSha512Mechanism(),
				),
			)
		case "SCRAM-SHA-256":
			opts = append(opts,
				kgo.SASL(
					scram.Auth{
						User: cfg.SASLUsername,
						Pass: cfg.SASLPassword,
					}.AsSha256Mechanism(),
				),
			)
		}
	}

	cl, err := kgo.NewClient(opts...)

	if err != nil {
		return nil, err
	}
	return &Consumer{
		cli: cl,
	}, nil
}

func (c *Consumer) Consume(
	ctx context.Context,
	iterator func(record *kgo.Record) error,
) error {
	for {
		select {
		case <-ctx.Done():
			c.cli.Close()
			return ctx.Err()
		default:
			fetches := c.cli.PollFetches(ctx)
			if fetches.IsClientClosed() {
				return fmt.Errorf("client closed")
			}
			if errs := fetches.Errors(); len(errs) > 0 {
				for _, err := range errs {
					return err.Err
				}
			}
			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()
				err := iterator(record)
				if err != nil {
					return err
				}
			}
		}
	}
}
