// Code generated from Pkl module `com.pipelaner.source.inputs`. DO NOT EDIT.
package commitstrategy

import (
	"encoding"
	"fmt"
)

type CommitStrategy string

const (
	OneByOne      CommitStrategy = "one-by-one"
	MarkOnSuccess CommitStrategy = "mark-on-success"
	AutoCommit    CommitStrategy = "auto-commit"
)

// String returns the string representation of CommitStrategy
func (rcv CommitStrategy) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(CommitStrategy)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for CommitStrategy.
func (rcv *CommitStrategy) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "one-by-one":
		*rcv = OneByOne
	case "mark-on-success":
		*rcv = MarkOnSuccess
	case "auto-commit":
		*rcv = AutoCommit
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid CommitStrategy`, str)
	}
	return nil
}
