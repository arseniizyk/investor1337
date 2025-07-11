package bot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/arseniizyk/investor1337/internal/aggregator"
)

func format(res []aggregator.MarketInfo) string {
	var result strings.Builder

	sort.Slice(res, func(i, j int) bool {
		isEmptyI := len(res[i].Orders) == 0
		isEmptyJ := len(res[j].Orders) == 0

		if isEmptyI && isEmptyJ {
			return false
		}
		if isEmptyI {
			return false
		}
		if isEmptyJ {
			return true
		}

		return res[i].Orders[0].Price < res[j].Orders[0].Price
	})

	for _, output := range res {
		if len(output.Orders) == 0 {
			result.WriteString(fmt.Sprintf("%s не найдено предложений\n\n", output.Market))
			continue
		}

		result.WriteString(fmt.Sprintf("%s\n", output.Market))
		for _, pair := range output.Orders {
			result.WriteString(fmt.Sprintf("Price: $%.2f | %d\n", pair.Price, pair.Quantity))
		}
		result.WriteString("\n")
	}

	return result.String()
}
