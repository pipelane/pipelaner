// Code generated from Pkl module `pipelaner.source.sinks`. DO NOT EDIT.
package sink

import "github.com/apple/pkl-go/pkl"

func init() {
	pkl.RegisterMapping("pipelaner.source.sinks", Sinks{})
	pkl.RegisterMapping("pipelaner.source.sinks#ExampleConsole", ExampleConsoleImpl{})
	pkl.RegisterMapping("pipelaner.source.sinks#Console", ConsoleImpl{})
	pkl.RegisterMapping("pipelaner.source.sinks#Pipelaner", PipelanerImpl{})
	pkl.RegisterMapping("pipelaner.source.sinks#KafkaProducer", KafkaProducerImpl{})
	pkl.RegisterMapping("pipelaner.source.sinks#Clickhouse", ClickhouseImpl{})
}
