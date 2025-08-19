/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/common/saslmechanism"
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/input/commitstrategy"
	"github.com/pipelane/pipelaner/gen/source/input/isolationlevel"
	"github.com/pipelane/pipelaner/gen/source/input/strategy"
	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/plugin/kzerolog"
)

type Consumer struct {
	cli    *kgo.Client
	logger *zerolog.Logger
}

func NewConsumer(
	cfg input.Kafka,
	logger *zerolog.Logger,
) (*Consumer, error) {
	maxPartitionFetchBytes := cfg.GetMaxPartitionFetchBytes()
	v := maxPartitionFetchBytes.ToUnit(pkl.Bytes).Value
	fetchMaxBytes := cfg.GetFetchMaxBytes()
	maxByteFetch := fetchMaxBytes.ToUnit(pkl.Bytes).Value
	isolationLevel := kgo.ReadUncommitted()
	if cfg.GetIsolationLevel() == isolationlevel.ReadCommitted {
		isolationLevel = kgo.ReadCommitted()
	}
	cons := &Consumer{
		logger: logger,
	}
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.GetCommon().Brokers...),
		kgo.WithLogger(kzerolog.New(logger)),
		kgo.ConsumerGroup(cfg.GetConsumerGroupID()),
		kgo.ConsumeTopics(cfg.GetCommon().Topics...),
		kgo.FetchMaxBytes(int32(maxByteFetch)),
		kgo.FetchMaxPartitionBytes(int32(v)),
		kgo.HeartbeatInterval(time.Second),
		kgo.FetchIsolationLevel(isolationLevel),
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

	switch cfg.GetCommitSrategy().Strategy {
	case commitstrategy.MarkOnSuccess:
		duration := cfg.GetCommitSrategy().Interval
		interval := duration.GoDuration()
		opts = append(opts, kgo.AutoCommitMarks())
		opts = append(opts, kgo.AutoCommitInterval(interval))
		opts = append(opts, kgo.OnPartitionsRevoked(cons.revoked))
	case commitstrategy.AutoCommit:
		interval := cfg.GetCommitSrategy().Interval
		opts = append(opts, kgo.AutoCommitInterval(interval.GoDuration()))
	case commitstrategy.OneByOne:
		opts = append(opts, kgo.DisableAutoCommit())
	}

	if cfg.GetAutoOffsetReset() == "earliest" || cfg.GetAutoOffsetReset() == "" {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()))
	} else if cfg.GetAutoOffsetReset() == "latest" {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()))
	}

	if cfg.GetCommon().SaslAuth != nil {
		auth := scram.Auth{
			User: cfg.GetCommon().SaslAuth.SaslUsername,
			Pass: cfg.GetCommon().SaslAuth.SaslPassword,
		}
		var authOpt kgo.Opt
		switch cfg.GetCommon().SaslAuth.SaslMechanism {
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
	cons.cli = cl
	return cons, nil
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

func (c *Consumer) CommitRecords(ctx context.Context, rec *kgo.Record) error {
	return c.cli.CommitRecords(ctx, rec)
}

func (c *Consumer) MarkCommitRecords(rec *kgo.Record) {
	c.cli.MarkCommitRecords(rec)
}

func (c *Consumer) CommitMarkedOffsets(ctx context.Context) error {
	return c.cli.CommitMarkedOffsets(ctx)
}

func (c *Consumer) revoked(ctx context.Context, cl *kgo.Client, _ map[string][]int32) {
	if err := cl.CommitMarkedOffsets(ctx); err != nil {
		c.logger.Error().Err(err).Msg("failed to revoke marked offsets")
	}
}
