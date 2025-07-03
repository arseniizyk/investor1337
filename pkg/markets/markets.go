package markets

const MaxOutputs = 3

type Market interface {
	FindByHashName(string) (map[float64]int, error)
}
