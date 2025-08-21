/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package http

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/gen/source/sink/method"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/node"
	"github.com/pipelane/pipelaner/pipeline/source"
)

type Client struct {
	components.Logger
	cfg sink.Http
	cli *resty.Client
}

func init() {
	source.RegisterSink("http", &Client{})
}

func (c *Client) Init(cfg sink.Sink) error {
	httpCfg, ok := cfg.(sink.Http)
	if !ok {
		return fmt.Errorf("invalid http client config %T", cfg)
	}
	c.cfg = httpCfg
	c.cli = resty.New()
	return nil
}

func (c *Client) Sink(val any) error {
	var data any
	switch v := val.(type) {
	case node.AtomicMessage:
		err := c.Sink(v.Data())
		if err != nil {
			v.Error() <- v
			return err
		}
		v.Success() <- v
		return nil
	default:
		data = val
	}
	r := c.cli.R()
	switch c.cfg.GetMethod() {
	case method.POST, method.PUT, method.PATCH, method.DELETE:
		r.SetBody(data)
	}
	if c.cfg.GetHeaders() != nil {
		r = r.SetHeaders(*c.cfg.GetHeaders())
	}
	m := c.method(r)
	resp, err := m(c.cfg.GetUrl())
	if err != nil || resp.IsError() {
		c.Log().
			Err(err).
			Bool("is_error", resp.IsError()).
			Str("url", c.cfg.GetUrl()).
			Msg("failed to send http request")
		return err
	}
	c.Log().
		Debug().
		Str("url", c.cfg.GetUrl()).
		Str("status_code", resp.Status()).
		Str("body", string(resp.Body())).
		Msg("received http request")
	return nil
}

func (c *Client) method(r *resty.Request) func(url string) (*resty.Response, error) {
	switch c.cfg.GetMethod() {
	case method.POST:
		return r.Post
	case method.PUT:
		return r.Put
	case method.PATCH:
		return r.Patch
	case method.DELETE:
		return r.Delete
	case method.GET:
		return r.Get
	}
	return nil
}
