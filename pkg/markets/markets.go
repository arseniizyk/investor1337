package markets

import "context"

const MaxOutputs = 4

type Market interface {
	FindByHashName(ctx context.Context, name string) (map[float64]int, error)
}
