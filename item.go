/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import "context"

type Context struct {
	ctx      context.Context
	lateItem *LaneItem
}

func NewContext(ctx context.Context, lateItem *LaneItem) *Context {
	return &Context{ctx: ctx, lateItem: lateItem}
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) LaneItem() *LaneItem {
	return c.lateItem
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
	outputCh := l.runLoop.createOutput(output.cfg.BufferSize)
	output.runLoop.setInputChannel(outputCh)
	l.addOutputs(output)
}

func NewLaneItem(
	config *BaseLaneConfig,
) *LaneItem {
	return &LaneItem{
		runLoop: newRunLoop(config.BufferSize, config.Threads),
		cfg:     config,
	}
}
