package pipelane

import "context"

type Agent struct {
	tree *TreeLanes
	ctx  context.Context
}

func NewAgent(
	ctx context.Context,
	dataSource DataSource,
	file string,
) (*Agent, error) {
	t, err := NewTreeFrom(
		ctx,
		dataSource,
		file,
	)
	if err != nil {
		return nil, err
	}
	return &Agent{
		tree: t,
		ctx:  ctx,
	}, err
}

func (a *Agent) Serve() {
	for {
		select {
		case <-a.ctx.Done():
			break
		default:
		}
	}
}
