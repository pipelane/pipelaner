/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package tests

/*
type IntStruct struct {
	Counter uint64
}

type IntGenerator struct {
	inc uint64
}

func (i *IntGenerator) Init(_ *pipelaner.Context) error {
	i.inc = 0
	return nil
}
func (i *IntGenerator) Generate(ctx *pipelaner.Context, input chan<- any) {
	for {
		select {
		case <-ctx.Context().Done():
			break
		default:
			input <- &IntStruct{
				Counter: i.inc,
			}
		}
		i.inc++
		if i.inc%10 == 0 {
			time.Sleep(60 * time.Second)
		}
	}
}

type IntTwoGenerator struct {
}

func (i *IntTwoGenerator) Init(_ *pipelaner.Context) error {
	return nil
}

func (i *IntTwoGenerator) Generate(ctx *pipelaner.Context, input chan<- any) {
	for {
		select {
		case <-ctx.Context().Done():
			break
		default:
			input <- 2
		}
	}
}

type IntTransform struct {
}

func (i *IntTransform) New() pipelaner.Map {
	return &IntTransform{}
}

func (i *IntTransform) Map(_ *pipelaner.Context, val any) any {
	time.Sleep(time.Second)
	return val.(uint64) + 2
}

func (i *IntTransform) Init(_ *pipelaner.Context) error {
	return nil
}

type IntTransformEmpty struct {
}

func (i *IntTransformEmpty) New() pipelaner.Map {
	return &IntTransformEmpty{}
}

func (i *IntTransformEmpty) Map(_ *pipelaner.Context, val any) any {
	time.Sleep(time.Second)
	return val
}

func (i *IntTransformEmpty) Init(_ *pipelaner.Context) error {
	return nil
}*/
