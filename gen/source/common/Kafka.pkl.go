// Code generated from Pkl module `com.pipelaner.source.common`. DO NOT EDIT.
package common

import "github.com/pipelane/pipelaner/gen/source/common/saslmechanism"

type Kafka struct {
	SaslEnabled bool `pkl:"saslEnabled"`

	SaslMechanism saslmechanism.SASLMechanism `pkl:"saslMechanism"`

	SaslUsername *string `pkl:"saslUsername"`

	SaslPassword *string `pkl:"saslPassword"`

	Brokers string `pkl:"brokers"`

	Version *string `pkl:"version"`

	Topics []string `pkl:"topics"`

	SchemaRegistry string `pkl:"schemaRegistry"`
}
