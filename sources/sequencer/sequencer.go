package sequencer

import (
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/source"
)

func init() {
	source.RegisterSequencer("sequencer", &Sequencer{})
}

type Sequencer struct {
	components.Logger
}

func (s *Sequencer) Init(_ transform.Sequencer) error {
	return nil
}
