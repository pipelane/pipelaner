// Code generated from Pkl module `pipelaner.source.Common`. DO NOT EDIT.
package saslmechanism

import (
	"encoding"
	"fmt"
)

type SASLMechanism string

const (
	SCRAMSHA512 SASLMechanism = "SCRAM-SHA-512"
	SCRAMSHA256 SASLMechanism = "SCRAM-SHA-256"
	PLAIN       SASLMechanism = "PLAIN"
)

// String returns the string representation of SASLMechanism
func (rcv SASLMechanism) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(SASLMechanism)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for SASLMechanism.
func (rcv *SASLMechanism) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "SCRAM-SHA-512":
		*rcv = SCRAMSHA512
	case "SCRAM-SHA-256":
		*rcv = SCRAMSHA256
	case "PLAIN":
		*rcv = PLAIN
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid SASLMechanism`, str)
	}
	return nil
}
