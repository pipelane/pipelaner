/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package generator

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"

	"github.com/pipelane/pipelaner"
)

type ExecCfg struct {
	Exec []string `pipelane:"exec"`
}

type Exec struct {
	cfg    *pipelaner.BaseLaneConfig
	logger zerolog.Logger
}

func (c *Exec) Init(cfg *pipelaner.BaseLaneConfig) error {
	c.cfg = cfg
	c.logger = pipelaner.NewLogger()
	v := &ExecCfg{}
	err := cfg.ParseExtended(v)
	if err != nil {
		return err
	}
	return nil
}

func (c *Exec) Generate(ctx context.Context, input chan<- any) {
	var args []string
	cfg := c.cfg.Extended.(*ExecCfg)
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

	if err = cmd.Wait(); err != nil {
		c.logger.Error().Err(err).Msg("Exec: command wait error")
	}
}

func (c *Exec) readPipe(ctx context.Context, pipe io.Reader, input chan<- any) {
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
