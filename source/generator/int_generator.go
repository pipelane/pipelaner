package generator

import "context"

type IntGenerator struct {
	inc int
}

func (i IntGenerator) Generate(ctx context.Context) any {
	return 1
}
