package pipelane

type Generators map[string]Generator
type Transformers map[string]Transformer
type Sinks map[string]Sink

type DataSource struct {
	Generators Generators
	Transforms Transformers
	Sinks      Sinks
}
