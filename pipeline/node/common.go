package node

type Type string

type nodeCfg struct {
	enableMetrics bool
	callGC        bool
}

type Option func(*nodeCfg)

func WithMetrics() Option {
	return func(cfg *nodeCfg) {
		cfg.enableMetrics = true
	}
}

func WithCallGC() Option {
	return func(cfg *nodeCfg) {
		cfg.callGC = true
	}
}

func buildOptions(opts ...Option) *nodeCfg {
	cfg := &nodeCfg{}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
