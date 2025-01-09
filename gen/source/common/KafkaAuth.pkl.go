// Code generated from Pkl module `com.pipelaner.source.common`. DO NOT EDIT.
package common

import "github.com/pipelane/pipelaner/gen/source/common/saslmechanism"

type KafkaAuth struct {
	SaslMechanism saslmechanism.SASLMechanism `pkl:"saslMechanism"`

	SaslUsername string `pkl:"saslUsername"`

	SaslPassword string `pkl:"saslPassword"`
}
