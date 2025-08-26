/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/node"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

func init() {
	source.RegisterSink("kafka", &Kafka{})
}

type producer interface {
	Produce(
		ctx context.Context,
		r *kgo.Record,
		promise func(*kgo.Record, error),
	)
	BeginTransaction() error
	Flush(ctx context.Context) error
	EndTransaction(ctx context.Context, commit kgo.TransactionEndTry) error
	AbortBufferedRecords(ctx context.Context) error
	ProduceSync(ctx context.Context, rs ...*kgo.Record) kgo.ProduceResults
}

type Kafka struct {
	components.Logger
	cfg  sink.Kafka
	prod producer
}

func (k *Kafka) Init(cfg sink.Sink) error {
	kafkaCfg, ok := cfg.(sink.Kafka)
	if !ok {
		return fmt.Errorf("invalid kafka-producer config %T", cfg)
	}
	kafkaLogger := k.Log().With().Logger()
	p, err := NewProducer(kafkaCfg, &kafkaLogger)
	if err != nil {
		return fmt.Errorf("init kafka producer: %w", err)
	}
	k.prod = p
	k.cfg = kafkaCfg
	return nil
}

func (k *Kafka) write(ctx context.Context, message []byte) {
	for _, topic := range k.cfg.GetCommon().Topics {
		k.prod.Produce(ctx, &kgo.Record{
			Value: message,
			Topic: topic,
		}, func(record *kgo.Record, err error) {
			if err != nil {
				k.Log().Error().Err(err).Msg("failed to produce message")
				k.write(ctx, record.Value)
			}
		})
	}
}

func (k *Kafka) writeSync(ctx context.Context, message any) error {
	var messageBytes []byte
	var err error
	switch byteVal := message.(type) {
	case []byte:
		messageBytes = byteVal
	case string:
		messageBytes = []byte(byteVal)
	default:
		messageBytes, err = json.Marshal(message)
		if err != nil {
			k.Log().Error().Err(err).Msg("marshal val")
			return err
		}
	}
	for _, topic := range k.cfg.GetCommon().Topics {
		err = k.prod.ProduceSync(ctx, &kgo.Record{
			Value: messageBytes,
			Topic: topic,
		}).FirstErr()
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Kafka) transactionalWrite(ctx context.Context, message any) error {
	if err := k.prod.BeginTransaction(); err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	var finishChan chan node.AtomicData
	switch msg := message.(type) {
	case chan node.AtomicData:
		finishChan = make(chan node.AtomicData, cap(msg))
		var err error
		for v := range msg {
			err = k.writeSync(ctx, v)
			if err != nil {
				e := k.rollback(ctx)
				if e != nil {
					k.Log().Error().Err(e).Msg("error rolling back")
				}
				finishChan <- v
				break
			}
			finishChan <- v
		}
		close(finishChan)
		if err != nil {
			k.sendAtomicError(msg)
			k.sendAtomicError(finishChan)
			return err
		}
	case chan []byte:
		for v := range msg {
			err := k.writeSync(ctx, v)
			if err != nil {
				e := k.rollback(ctx)
				if e != nil {
					k.Log().Error().Err(e).Msg("error rolling back")
				}
				return err
			}
		}
	case chan string:
		for v := range msg {
			err := k.writeSync(ctx, v)
			if err != nil {
				e := k.rollback(ctx)
				if e != nil {
					k.Log().Error().Err(e).Msg("error rolling back")
				}
				return err
			}
		}
	case chan any:
		for v := range msg {
			err := k.writeSync(ctx, v)
			if err != nil {
				e := k.rollback(ctx)
				if e != nil {
					k.Log().Error().Err(e).Msg("error rolling back")
				}
				return err
			}
		}
	default:
		err := k.rollback(ctx)
		if err != nil {
			return err
		}
		return fmt.Errorf("invalid message type %T", message)
	}
	if err := k.prod.Flush(ctx); err != nil {
		return fmt.Errorf("kafka producer failed to commit transaction: %w", err)
	}
	err := k.prod.EndTransaction(ctx, kgo.TryCommit)
	if err != nil && errors.Is(err, kerr.OperationNotAttempted) {
		err = k.rollback(ctx)
		if err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	k.sendAtomicSuccess(finishChan)
	return nil
}

func (k *Kafka) Sink(val any) error {
	var message []byte
	var err error

	switch v := val.(type) {
	case node.AtomicData:
		err = k.Sink(v.Data())
		if err != nil {
			v.Error() <- v
			return err
		}
		v.Success() <- v
		return nil
	case []byte:
		message = v
	case string:
		message = []byte(v)
	case chan node.AtomicData:
		err = k.transactionalWrite(context.Background(), v)
		if err != nil {
			return err
		}
		return nil
	case chan []byte:
		err = k.transactionalWrite(context.Background(), v)
		if err != nil {
			return err
		}
		return nil
	case chan string:
		err = k.transactionalWrite(context.Background(), v)
		if err != nil {
			return err
		}
		return nil
	case chan any:
		err = k.transactionalWrite(context.Background(), v)
		if err != nil {
			return err
		}
		return nil
	default:
		data, errs := json.Marshal(val)
		if errs != nil {
			k.Log().Error().Err(errs).Msg("marshal val")
			return errs
		}
		message = data
	}
	k.write(context.Background(), message)
	return nil
}

func (k *Kafka) rollback(ctx context.Context) error {
	if err := k.prod.AbortBufferedRecords(ctx); err != nil {
		return fmt.Errorf("error rolling back buffered records: %w", err)
	}
	if err := k.prod.EndTransaction(ctx, kgo.TryAbort); err != nil {
		return fmt.Errorf("error rolling back transaction: %w", err)
	}
	return nil
}

func (k *Kafka) sendAtomicError(chData chan node.AtomicData) {
	for chV := range chData {
		switch vals := chV.(type) {
		case node.AtomicData:
			vals.Error() <- vals
		default:
			break
		}
	}
}

func (k *Kafka) sendAtomicSuccess(chData chan node.AtomicData) {
	if chData == nil {
		return
	}
	for chV := range chData {
		switch vals := chV.(type) {
		case node.AtomicData:
			vals.Success() <- vals
		default:
			break
		}
	}
}
