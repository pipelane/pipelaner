// Code generated from Pkl module `pipelaner.source.example`. DO NOT EDIT.
package custom

import "github.com/apple/pkl-go/pkl"

func init() {
	pkl.RegisterMapping("pipelaner.source.example", Example{})
	pkl.RegisterMapping("pipelaner.source.example#ExampleGenInt", ExampleGenIntImpl{})
	pkl.RegisterMapping("pipelaner.source.example#ExampleMul", ExampleMulImpl{})
	pkl.RegisterMapping("pipelaner.source.example#ExampleConsole", ExampleConsoleImpl{})
}
