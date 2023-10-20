package generator

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"

	pipelane "github.com/pipelane/pipelaner"
)

type CommandCfg struct {
	Exec []string `pipelane:"exec"`
}

type Command struct {
	cfg    *pipelane.BaseLaneConfig
	logger zerolog.Logger
}

func (c *Command) Init(cfg *pipelane.BaseLaneConfig) error {
	c.cfg = cfg
	c.logger = pipelane.NewLogger()
	v := &CommandCfg{}
	err := cfg.ParseExtended(v)
	if err != nil {
		return err
	}
	return nil
}

func (c *Command) Generate(ctx context.Context, input chan<- any) {
	var args []string
	cfg := c.cfg.Extended.(*CommandCfg)
	if len(cfg.Exec) > 1 {
		args = strings.Split(cfg.Exec[1], " ")
	}
	cmd := exec.Command(cfg.Exec[0], args...) //nolint:gosec
	stdPipe, err := cmd.StdoutPipe()
	if err != nil {
		c.logger.Error().Err(err).Msg("Exec: create stdPipe error")
		return
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		c.logger.Error().Err(err).Msg("Exec: create stdPipe error")
		return
	}
	if err = cmd.Start(); err != nil {
		c.logger.Error().Err(err).Msg("Exec: create errPipe error")
		return
	}
	go c.readPipe(ctx, stdPipe, input)
	go c.readPipe(ctx, stdErr, input)

	if err := cmd.Wait(); err != nil {
		c.logger.Error().Err(err).Msg("Exec: command wait error")
	}
}

func (c *Command) readPipe(ctx context.Context, pipe io.Reader, input chan<- any) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		select {
		case <-ctx.Done():
			break
		default:
			input <- m
		}
	}
}
