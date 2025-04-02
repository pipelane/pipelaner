// Code generated from Pkl module `com.pipelaner.source.transforms`. DO NOT EDIT.
package transform

import "github.com/apple/pkl-go/pkl"

func init() {
	pkl.RegisterMapping("com.pipelaner.source.transforms", Transforms{})
	pkl.RegisterMapping("com.pipelaner.source.transforms#Sequencer", SequencerImpl{})
	pkl.RegisterMapping("com.pipelaner.source.transforms#Batch", BatchImpl{})
	pkl.RegisterMapping("com.pipelaner.source.transforms#Chunk", ChunkImpl{})
	pkl.RegisterMapping("com.pipelaner.source.transforms#Debounce", DebounceImpl{})
	pkl.RegisterMapping("com.pipelaner.source.transforms#Throttling", ThrottlingImpl{})
	pkl.RegisterMapping("com.pipelaner.source.transforms#Filter", FilterImpl{})
	pkl.RegisterMapping("com.pipelaner.source.transforms#Remap", RemapImpl{})
}
