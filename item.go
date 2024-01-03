/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

type LaneItem struct {
	runLoop         *runLoop
	Cfg             *BaseLaneConfig
	inputPipeline   []*LaneItem
	outputPipelines []*LaneItem
}

func (p *LaneItem) addInputPipelines(i *LaneItem) {
	p.inputPipeline = append(p.inputPipeline, i)
}

func (p *LaneItem) addOutputs(output *LaneItem) {
	p.outputPipelines = append(p.outputPipelines, output)
	output.addInputPipelines(p)
}

func (p *LaneItem) Subscribe(output *LaneItem) {
	outputCh := p.runLoop.createOutput(output.Cfg.BufferSize)
	output.runLoop.setInputChannel(outputCh)
	p.addOutputs(output)
}

func NewLaneItem(
	config *BaseLaneConfig,
) *LaneItem {
	return &LaneItem{
		runLoop: newRunLoop(config.BufferSize, config.Threads),
		Cfg:     config,
	}
}
