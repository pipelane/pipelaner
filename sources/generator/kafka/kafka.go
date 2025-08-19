/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"

	inputsCfg "github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/input/commitstrategy"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/node"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/twmb/franz-go/pkg/kgo"
)

func init() {
	source.RegisterInput("kafka", &Kafka{})
}

type Kafka struct {
	components.Logger
	cons         *Consumer
	cfg          inputsCfg.Kafka
	consumeStore *sync.Map
}

func (c *Kafka) Init(cfg inputsCfg.Input) error {
	consumerCfg, ok := cfg.(inputsCfg.Kafka)
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
	var (
		successCh chan node.AtomicData
		errorsCh  chan node.AtomicData
	)
	l := c.Log()
	switch c.cfg.GetCommitSrategy().Strategy {
	case commitstrategy.OneByOne:
		successCh = make(chan node.AtomicData, 1)
		errorsCh = make(chan node.AtomicData, 1)
		defer close(successCh)
		defer close(errorsCh)
	case commitstrategy.MarkOnSuccess:
		successCh = make(chan node.AtomicData, c.cfg.GetOutputBufferSize())
		errorsCh = make(chan node.AtomicData, c.cfg.GetOutputBufferSize())
		c.consumeStore = &sync.Map{}
		go c.markRecord(successCh, errorsCh) //nolint: contextcheck
		defer close(successCh)
		defer close(errorsCh)
	case commitstrategy.AutoCommit:
	}
	for {
		err := c.cons.Consume(ctx, func(record *kgo.Record) error {
			switch c.cfg.GetCommitSrategy().Strategy {
			case commitstrategy.OneByOne:
				inputValue := node.NewAtomicMessage(record.Value, successCh, errorsCh)
				input <- inputValue
				select {
				case <-successCh:
					return c.cons.CommitRecords(ctx, record)
				case <-errorsCh:
					err := errors.New("failed processing message")
					c.Log().Error().Err(err).
						Int64("offset", record.Offset).
						Int32("partition", record.Partition).
						Msg("failed processing message")
					return err
				}
			case commitstrategy.AutoCommit:
				input <- record.Value
				return nil
			case commitstrategy.MarkOnSuccess:
				inputValue := node.NewAtomicMessage(record.Value, successCh, errorsCh)
				c.consumeStore.Store(inputValue.ID(), record)
				input <- inputValue
				return nil
			}
			return nil
		})
		c.commitMarked(context.Background()) //nolint: contextcheck
		if errors.Is(err, context.Canceled) {
			break
		}
		if err != nil {
			l.Error().Err(err).Msg("consume error")
		}
	}
}

func (c *Kafka) markRecord(successCh chan node.AtomicData, errCh chan node.AtomicData) {
Loop:
	for {
		select {
		case message, isClosed := <-successCh:
			if isClosed && message == nil {
				break Loop
			}
			val, ok := c.consumeStore.Load(message.ID())
			if !ok {
				panic("failed processing message")
			}
			v, ok := val.(*kgo.Record)
			if !ok {
				panic("failed processing message")
			}
			c.cons.MarkCommitRecords(v)
			c.consumeStore.Delete(message.ID())
		case message, isClosed := <-errCh:
			if isClosed && message == nil {
				break Loop
			}
			val, ok := c.consumeStore.Load(message.ID())
			if !ok {
				panic("failed processing message")
			}
			v, ok := val.(*kgo.Record)
			if !ok {
				panic("failed processing message")
			}
			err := errors.New("failed processing message")
			c.Log().Error().Err(err).
				Int64("offset", v.Offset).
				Int32("partition", v.Partition).
				Msg("failed processing message")
			c.consumeStore.Delete(message.ID())
		}
	}
	c.commitMarked(context.Background())
}

func (c *Kafka) commitMarked(ctx context.Context) {
	if c.cfg.GetCommitSrategy().Strategy == commitstrategy.MarkOnSuccess {
		err := c.cons.CommitMarkedOffsets(ctx)
		if err != nil {
			c.Log().Error().Err(err).Msg("failed commit marked messages")
		}
	}
}
