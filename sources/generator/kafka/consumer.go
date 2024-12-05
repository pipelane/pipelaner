package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/common/saslmechanism"
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/input/strategy"
	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/plugin/kzerolog"
)

type Consumer struct {
	cli *kgo.Client
}

func NewConsumer(
	cfg input.KafkaConsumer,
	logger *zerolog.Logger,
) (*Consumer, error) {
	v := cfg.GetMaxPartitionFetchBytes().ToUnit(pkl.Bytes).Value
	maxByteFetch := cfg.GetFetchMaxBytes().ToUnit(pkl.Bytes).Value

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.GetKafka().Brokers),
		kgo.WithLogger(kzerolog.New(logger)),
		kgo.ConsumerGroup(cfg.GetConsumerGroupID()),
		kgo.ConsumeTopics(cfg.GetKafka().Topics...),
		kgo.FetchMaxBytes(int32(maxByteFetch)),
		kgo.FetchMaxPartitionBytes(int32(v)),
		kgo.HeartbeatInterval(time.Second),
	}
	var balancers []kgo.GroupBalancer
	for _, s := range cfg.GetBalancerStrategy() {
		switch s {
		case strategy.Range:
			balancers = append(balancers, kgo.RangeBalancer())
		case strategy.RoundRobin:
			balancers = append(balancers, kgo.RoundRobinBalancer())
		case strategy.CooperativeSticky:
			balancers = append(balancers, kgo.CooperativeStickyBalancer())
		case strategy.Sticky:
			balancers = append(balancers, kgo.StickyBalancer())
		}
	}
	if len(balancers) == 0 {
		balancers = append(balancers, kgo.CooperativeStickyBalancer())
	}
	opts = append(opts, kgo.Balancers(balancers...))

	if !cfg.GetAutoCommitEnabled() {
		opts = append(opts, kgo.DisableAutoCommit())
	}

	if cfg.GetAutoOffsetReset() == "earliest" || cfg.GetAutoOffsetReset() == "" {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()))
	} else if cfg.GetAutoOffsetReset() == "latest" {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()))
	}

	if cfg.GetKafka().SaslEnabled {
		auth := scram.Auth{
			User: *cfg.GetKafka().SaslUsername,
			Pass: *cfg.GetKafka().SaslPassword,
		}
		var authOpt kgo.Opt
		switch cfg.GetKafka().SaslMechanism {
		case saslmechanism.SCRAMSHA512:
			authOpt = kgo.SASL(auth.AsSha512Mechanism())
		case saslmechanism.SCRAMSHA256:
			authOpt = kgo.SASL(auth.AsSha256Mechanism())
		}
		opts = append(opts, authOpt)
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
