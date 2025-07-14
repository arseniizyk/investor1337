package markets

import (
	"sort"
)

type Pair struct {
	Price    float64
	Quantity int
}

func SinglePair(price float64, count int) []Pair {
	return []Pair{{
		Price:    price,
		Quantity: count,
	}}
}

func SortPairs(pairs []Pair) {
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].Price < pairs[j].Price
	})
}

func PairsFromMap(m map[float64]int) []Pair {
	result := make([]Pair, 0, len(m))

	for price, quantity := range m {
		result = append(result, Pair{
			Price:    price,
			Quantity: quantity,
		})
	}

	SortPairs(result)

	return result
}
