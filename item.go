/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
)

var (
	ErrLaneIsStopped = errors.New("ErrLaneIsStopped")
)

type Context struct {
	ctx      context.Context
	laneItem *LaneItem
	value    any
	cancel   context.CancelFunc
}

func NewContext(ctx context.Context, laneItem *LaneItem) *Context {
	c, cancel := context.WithCancel(ctx)
	return &Context{ctx: c, laneItem: laneItem, cancel: cancel}
}
func (c *Context) LaneName() string {
	return c.laneItem.cfg.Name
}

func (c *Context) LaneType() LaneTypes {
	return c.laneItem.cfg.LaneType
}
func (c *Context) SourceName() string {
	return c.laneItem.cfg.SourceName
}

func withContext(parent *Context) *Context {
	c, cancel := context.WithCancel(parent.ctx)
	return &Context{ctx: c, cancel: cancel}
}

func (c *Context) Value() any {
	return c.value
}

func (c *Context) Logger() zerolog.Logger {
	v, ok := c.ctx.Value(logKey).(*zerolog.Logger)
	if !ok {
		return zerolog.Nop()
	}
	ctx := v.With().
		Str("source", c.SourceName()).
		Str("lane_name", c.LaneName()).
		Str("type", string(c.LaneType()))
	return ctx.Logger()
}

func (c *Context) ReturnValue(value any) error {
	if !c.LaneItem().runLoop.stopped.Load() {
		c.value = value
	}
	return ErrLaneIsStopped
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) LaneItem() *LaneItem {
	return c.laneItem
}

type LaneItem struct {
	runLoop         *runLoop
	cfg             *BaseLaneConfig
	inputPipeline   []*LaneItem
	outputPipelines []*LaneItem
}

func (l *LaneItem) Config() *BaseLaneConfig {
	return l.cfg
}

func (l *LaneItem) addInputPipelines(i *LaneItem) {
	l.inputPipeline = append(l.inputPipeline, i)
}

func (l *LaneItem) addOutputs(output *LaneItem) {
	l.outputPipelines = append(l.outputPipelines, output)
	output.addInputPipelines(l)
}

func (l *LaneItem) Subscribe(output *LaneItem) {
	outputCh := l.runLoop.createOutput(output.cfg.OutputBufferSize)
	output.runLoop.setInputChannel(outputCh)
	ctx := withContext(l.runLoop.context)
	ctx.laneItem = output
	output.runLoop.setContext(ctx)
}

func NewLaneItem(
	config *BaseLaneConfig,
) *LaneItem {
	return &LaneItem{
		runLoop: newRunLoop(config.OutputBufferSize, config.Threads, config.StartGCAfterMessageProcess),
		cfg:     config,
	}
}
