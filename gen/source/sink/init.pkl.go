// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package sink

import "github.com/apple/pkl-go/pkl"

func init() {
	pkl.RegisterMapping("com.pipelaner.source.sinks", Sinks{})
	pkl.RegisterMapping("com.pipelaner.source.sinks#Console", ConsoleImpl{})
	pkl.RegisterMapping("com.pipelaner.source.sinks#Pipelaner", PipelanerImpl{})
	pkl.RegisterMapping("com.pipelaner.source.sinks#KafkaProducer", KafkaProducerImpl{})
	pkl.RegisterMapping("com.pipelaner.source.sinks#Clickhouse", ClickhouseImpl{})
}
