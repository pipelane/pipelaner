// Code generated from Pkl module `com.pipelaner.source.inputs`. DO NOT EDIT.
package strategy

import (
	"encoding"
	"fmt"
)

type Strategy string

const (
	Range             Strategy = "range"
	RoundRobin        Strategy = "round-robin"
	CooperativeSticky Strategy = "cooperative-sticky"
	Sticky            Strategy = "sticky"
)

// String returns the string representation of Strategy
func (rcv Strategy) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(Strategy)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for Strategy.
func (rcv *Strategy) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "range":
		*rcv = Range
	case "round-robin":
		*rcv = RoundRobin
	case "cooperative-sticky":
		*rcv = CooperativeSticky
	case "sticky":
		*rcv = Sticky
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid Strategy`, str)
	}
	return nil
}
