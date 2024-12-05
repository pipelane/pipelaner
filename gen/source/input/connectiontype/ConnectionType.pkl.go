// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package connectiontype

import (
	"encoding"
	"fmt"
)

type ConnectionType string

const (
	Unix  ConnectionType = "unix"
	Http2 ConnectionType = "http2"
)

// String returns the string representation of ConnectionType
func (rcv ConnectionType) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(ConnectionType)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for ConnectionType.
func (rcv *ConnectionType) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "unix":
		*rcv = Unix
	case "http2":
		*rcv = Http2
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid ConnectionType`, str)
	}
	return nil
}
