// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package input

type KafkaConfig struct {
	SaslEnabled *bool `pkl:"saslEnabled"`

	SaslMechanism *bool `pkl:"saslMechanism"`

	SaslUsername *string `pkl:"saslUsername"`

	SaslPassword *string `pkl:"saslPassword"`

	Brokers string `pkl:"brokers"`

	Version string `pkl:"version"`

	Topics []string `pkl:"topics"`

	SchemaRegistry string `pkl:"schemaRegistry"`
}
