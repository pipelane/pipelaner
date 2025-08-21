// Code generated from Pkl module `com.pipelaner.source.inputs`. DO NOT EDIT.
package input

import (
	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/input/commitstrategy"
)

type CommitConfiguration struct {
	Strategy commitstrategy.CommitStrategy `pkl:"strategy"`

	Interval pkl.Duration `pkl:"interval"`
}
