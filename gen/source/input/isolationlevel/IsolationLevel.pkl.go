// Code generated from Pkl module `com.pipelaner.source.inputs`. DO NOT EDIT.
package isolationlevel

import (
	"encoding"
	"fmt"
)

type IsolationLevel string

const (
	ReadCommitted   IsolationLevel = "read-committed"
	ReadUncommitted IsolationLevel = "read-uncommitted"
)

// String returns the string representation of IsolationLevel
func (rcv IsolationLevel) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(IsolationLevel)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for IsolationLevel.
func (rcv *IsolationLevel) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "read-committed":
		*rcv = ReadCommitted
	case "read-uncommitted":
		*rcv = ReadUncommitted
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid IsolationLevel`, str)
	}
	return nil
}
