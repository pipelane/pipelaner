// Code generated from Pkl module `com.pipelaner.settings.logger.LoggerConfig`. DO NOT EDIT.
package logformat

import (
	"encoding"
	"fmt"
)

type LogFormat string

const (
	Plain LogFormat = "plain"
	Json  LogFormat = "json"
)

// String returns the string representation of LogFormat
func (rcv LogFormat) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(LogFormat)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for LogFormat.
func (rcv *LogFormat) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "plain":
		*rcv = Plain
	case "json":
		*rcv = Json
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid LogFormat`, str)
	}
	return nil
}
