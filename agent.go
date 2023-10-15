package pipelane

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Agent struct {
	tree   *TreeLanes
	ctx    context.Context
	cancel context.CancelFunc
}

func NewAgent(
	dataSource DataSource,
	file string,
) (*Agent, error) {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	t, err := NewTreeFrom(
		ctx,
		dataSource,
		file,
	)
	if err != nil {
		return nil, err
	}
	return &Agent{
		tree:   t,
		ctx:    ctx,
		cancel: stop,
	}, err
}

func (a *Agent) Serve() {
	<-a.ctx.Done()
	time.Sleep(time.Second * 10)
}

func (a *Agent) Stop() {
	a.cancel()
}
