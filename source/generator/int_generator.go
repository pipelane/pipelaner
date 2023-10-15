package generator

import "context"

type IntGenerator struct {
}

func (i IntGenerator) Generate(ctx context.Context) any {
	return 1
}
