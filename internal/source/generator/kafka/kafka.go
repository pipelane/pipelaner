package kafka

/*func init() {
	source.RegisterInput("kafka-consumer", &Kafka{})
}

type Kafka struct {
	cons *Consumer
	cfg  kcfg.ConsumerConfig
}

func (c *Kafka) Init(cfg input.Input) error {
	consumerCfg, err :=
}

func (c *Kafka) Init(ctx *pipelaner.Context) error {
	err := ctx.LaneItem().Config().ParseExtended(&c.cfg)
	if err != nil {
		return err
	}
	l := ctx.Logger()
	c.cons, err = NewConsumer(c.cfg, &l)
	if err != nil {
		return err
	}

	return nil
}

func (c *Kafka) Generate(ctx *pipelaner.Context, input chan<- any) {
	l := ctx.Logger()
	for {
		err := c.cons.Consume(ctx.Context(), func(record *kgo.Record) error {
			input <- record.Value
			return nil
		})
		if err != nil {
			l.Log().Err(err).Msg("consume error")
		}
	}
}*/