package markets

const MaxOutputs = 4

type Market interface {
	FindByHashName(string) (map[float64]int, error)
}
