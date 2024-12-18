// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package method

import (
	"encoding"
	"fmt"
)

type Method string

const (
	PATCH  Method = "PATCH"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	GET    Method = "GET"
)

// String returns the string representation of Method
func (rcv Method) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(Method)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for Method.
func (rcv *Method) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "PATCH":
		*rcv = PATCH
	case "POST":
		*rcv = POST
	case "PUT":
		*rcv = PUT
	case "DELETE":
		*rcv = DELETE
	case "GET":
		*rcv = GET
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid Method`, str)
	}
	return nil
}
