// Code generated from Pkl module `com.pipelaner.source.transforms`. DO NOT EDIT.
package flushtimerstrategy

import (
	"encoding"
	"fmt"
)

type FlushTimerStrategy string

const (
	FlushByOneMessage FlushTimerStrategy = "flush-by-one-message"
	FlushByTime       FlushTimerStrategy = "flush-by-time"
)

// String returns the string representation of FlushTimerStrategy
func (rcv FlushTimerStrategy) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(FlushTimerStrategy)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for FlushTimerStrategy.
func (rcv *FlushTimerStrategy) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "flush-by-one-message":
		*rcv = FlushByOneMessage
	case "flush-by-time":
		*rcv = FlushByTime
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid FlushTimerStrategy`, str)
	}
	return nil
}
