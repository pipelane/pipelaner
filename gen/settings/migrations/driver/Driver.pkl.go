// Code generated from Pkl module `com.pipelaner.settings.migrations.config`. DO NOT EDIT.
package driver

import (
	"encoding"
	"fmt"
)

type Driver string

const (
	Clickhouse Driver = "clickhouse"
	Empty      Driver = ""
)

// String returns the string representation of Driver
func (rcv Driver) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(Driver)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for Driver.
func (rcv *Driver) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "clickhouse":
		*rcv = Clickhouse
	case "":
		*rcv = Empty
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid Driver`, str)
	}
	return nil
}
