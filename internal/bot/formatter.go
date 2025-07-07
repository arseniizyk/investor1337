package bot

import (
	"fmt"
	"sort"
	"strings"
)

func format(res map[string]map[float64]int) string {
	var result strings.Builder

	for market, offers := range res {
		if len(offers) == 0 {
			result.WriteString(fmt.Sprintf("%s не найдено предложений\n\n", market))
			continue
		}

		keys := make([]float64, 0, len(offers))
		for k := range offers {
			keys = append(keys, k)
		}
		sort.Float64s(keys)

		result.WriteString(fmt.Sprintf("%s\n", market))
		for _, k := range keys {
			result.WriteString(fmt.Sprintf("Price: $%.2f | %d\n", k, offers[k]))
		}
		result.WriteString("\n")
	}

	return result.String()
}
