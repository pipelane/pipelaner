// Code generated from Pkl module `pipelaner.source.transforms`. DO NOT EDIT.
package transform

import "github.com/apple/pkl-go/pkl"

func init() {
	pkl.RegisterMapping("pipelaner.source.transforms", Transforms{})
	pkl.RegisterMapping("pipelaner.source.transforms#Batch", BatchImpl{})
	pkl.RegisterMapping("pipelaner.source.transforms#Chunk", ChunkImpl{})
	pkl.RegisterMapping("pipelaner.source.transforms#Debounce", DebounceImpl{})
	pkl.RegisterMapping("pipelaner.source.transforms#Throttling", ThrottlingImpl{})
	pkl.RegisterMapping("pipelaner.source.transforms#Filter", FilterImpl{})
	pkl.RegisterMapping("pipelaner.source.transforms#Remap", RemapImpl{})
	pkl.RegisterMapping("pipelaner.source.transforms#Mul", MulImpl{})
}
