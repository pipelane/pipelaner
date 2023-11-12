/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import "context"

type LaneItem struct {
	runLoop         *runLoop
	Cfg             *BaseLaneConfig
	inputPipeline   *LaneItem
	outputPipelines []*LaneItem
}

func (p *LaneItem) setInputPipelines(i *LaneItem) {
	p.inputPipeline = i
}

func (p *LaneItem) addOutputs(output *LaneItem) {
	p.outputPipelines = append(p.outputPipelines, output)
	output.setInputPipelines(p)
}

func (p *LaneItem) Subscribe(output *LaneItem) {
	outputCh := p.runLoop.createOutput(output.Cfg.BufferSize)
	output.runLoop.setInputChannel(outputCh)
	p.addOutputs(output)
}

func NewLaneItem(
	ctx context.Context,
	config *BaseLaneConfig,
) *LaneItem {
	return &LaneItem{
		runLoop: newRunLoop(ctx, config.BufferSize, config.Threads),
		Cfg:     config,
	}
}
