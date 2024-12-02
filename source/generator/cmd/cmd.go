/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/internal/logger"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/rs/zerolog"
)

func init() {
	source.RegisterInput("cmd", &Cmd{})
}

type Cmd struct {
	exec   []string
	logger *zerolog.Logger
}

func (c *Cmd) Init(cfg input.Input) error {
	cmdConfig, ok := cfg.(input.Cmd)
	if !ok {
		return fmt.Errorf("invalid cmd config type: %T", cfg)
	}
	// todo: inject logger here
	l := logger.NewLogger()
	c.logger = &l
	c.exec = cmdConfig.GetExec()
	return nil
}

func (c *Cmd) Generate(ctx context.Context, input chan<- any) {
	var args []string
	if len(c.exec) > 1 {
		args = strings.Split(c.exec[1], " ")
	}
	cmd := exec.Command(c.exec[0], args...) //nolint:gosec
	stdPipe, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		c.logger.Error().Err(err).Msg("Cmd: create stdPipe error")
		return
	}
	if err = cmd.Start(); err != nil {
		c.logger.Error().Err(err).Msg("Cmd: create errPipe error")
		return
	}
	go c.readPipe(ctx, stdPipe, input)
	go c.readPipe(ctx, stdErr, input)

	if err = cmd.Wait(); err != nil {
		c.logger.Error().Err(err).Msg("Cmd: command wait error")
	}
}

func (c *Cmd) readPipe(ctx context.Context, pipe io.Reader, input chan<- any) {
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
