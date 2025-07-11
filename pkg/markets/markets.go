package markets

import "context"

const MaxOutputs = 4

type Market interface {
	FindByHashName(ctx context.Context, name string) ([]Pair, error)
	URL(name string) string
}

type Pair struct {
	Price    float64
	Quantity int
}
