// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package sink

import "github.com/apple/pkl-go/pkl"

func init() {
	pkl.RegisterStrictMapping("com.pipelaner.source.sinks", Sinks{})
	pkl.RegisterStrictMapping("com.pipelaner.source.sinks#Console", ConsoleImpl{})
	pkl.RegisterStrictMapping("com.pipelaner.source.sinks#Pipelaner", PipelanerImpl{})
	pkl.RegisterStrictMapping("com.pipelaner.source.sinks#Kafka", KafkaImpl{})
	pkl.RegisterStrictMapping("com.pipelaner.source.sinks#Clickhouse", ClickhouseImpl{})
	pkl.RegisterStrictMapping("com.pipelaner.source.sinks#Http", HttpImpl{})
}
