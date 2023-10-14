package pipelane

import "context"

type LaneItem struct {
	*subscriber
	Cfg             *BaseConfig
	inputPipeline   *LaneItem
	outputPipelines []*LaneItem
}

func (p *LaneItem) setInputPipelines(i *LaneItem) {
	p.inputPipeline = i
}

func (p *LaneItem) outputs() []*LaneItem {
	return p.outputPipelines
}

func (p *LaneItem) addOutputs(output *LaneItem) {
	p.outputPipelines = append(p.outputPipelines, output)
	output.setInputPipelines(p)
}

func (p *LaneItem) Subscribe(output *LaneItem) {
	outputCh := p.createOutput(output.Cfg.BufferSize)
	output.setInputChannel(outputCh)
	p.addOutputs(output)
}

func NewLaneItem(
	ctx context.Context,
	config *BaseConfig,
) *LaneItem {
	return &LaneItem{
		subscriber: newSubscriber(ctx, config.BufferSize, config.ThreadsCount),
		Cfg:        config,
	}
}
