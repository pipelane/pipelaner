// Code generated from Pkl module `com.pipelaner.source.common`. DO NOT EDIT.
package common

type Kafka struct {
	SaslAuth *KafkaAuth `pkl:"saslAuth"`

	Brokers []string `pkl:"brokers"`

	Version *string `pkl:"version"`

	Topics []string `pkl:"topics"`
}
