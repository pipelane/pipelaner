/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package cmd

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os/exec"
	"strings"

	"github.com/pipelane/pipelaner"
)

type Config struct {
	Exec []string `pipelane:"exec"`
}

type Exec struct {
	cfg *pipelaner.BaseLaneConfig
}

func init() {
	pipelaner.RegisterGenerator("cmd", &Exec{})
}

func (c *Exec) Init(ctx *pipelaner.Context) error {
	c.cfg = ctx.LaneItem().Config()
	v := &Config{}
	err := c.cfg.ParseExtended(v)
	if err != nil {
		return err
	}
	return nil
}

func (c *Exec) Generate(ctx *pipelaner.Context, input chan<- any) {
	l := ctx.Logger()
	var args []string
	cfg, ok := c.cfg.Extended.(*Config)
	if !ok {
		l.Error().Err(errors.New("invalid config")).Msg("Exec: create stdPipe error")
		return
	}
	if len(cfg.Exec) > 1 {
		args = strings.Split(cfg.Exec[1], " ")
	}
	cmd := exec.Command(cfg.Exec[0], args...) //nolint:gosec
	stdPipe, err := cmd.StdoutPipe()
	if err != nil {
		l.Error().Err(err).Msg("Exec: create stdPipe error")
		return
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		l.Error().Err(err).Msg("Exec: create stdPipe error")
		return
	}
	if err = cmd.Start(); err != nil {
		l.Error().Err(err).Msg("Exec: create errPipe error")
		return
	}
	go c.readPipe(ctx.Context(), stdPipe, input)
	go c.readPipe(ctx.Context(), stdErr, input)

	if err = cmd.Wait(); err != nil {
		l.Error().Err(err).Msg("Exec: command wait error")
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
