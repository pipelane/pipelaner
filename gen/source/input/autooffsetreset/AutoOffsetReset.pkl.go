// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package autooffsetreset

import (
	"encoding"
	"fmt"
)

type AutoOffsetReset string

const (
	Earliest AutoOffsetReset = "earliest"
	Latest   AutoOffsetReset = "latest"
)

// String returns the string representation of AutoOffsetReset
func (rcv AutoOffsetReset) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(AutoOffsetReset)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for AutoOffsetReset.
func (rcv *AutoOffsetReset) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "earliest":
		*rcv = Earliest
	case "latest":
		*rcv = Latest
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid AutoOffsetReset`, str)
	}
	return nil
}
