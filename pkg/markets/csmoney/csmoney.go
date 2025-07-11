package csmoney

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

func (csm csmoney) FindByHashName(ctx context.Context, name string) ([]markets.Pair, error) {
	endpoint := "https://cs.money/2.0/market/sell-orders"
	params := url.Values{
		"limit":  []string{"60"},
		"offset": []string{"0"},
		"name":   []string{name},
		"order":  []string{"asc"},
		"sort":   []string{"price"},
	}

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	if err != nil {
		csm.l.Error("Cant make request to cs.money",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	r, err := u.DoJSONRequest[Response](ctx, csm.client, req, csm.l)

	if err != nil {
		csm.l.Warn("Response error from cs.money",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	return format(&r), nil
}

func format(r *Response) []markets.Pair {
	countMap := make(map[float64]int, markets.MaxOutputs)

	for _, item := range r.Items {
		if len(countMap) == markets.MaxOutputs {
			break
		}
		countMap[item.Pricing.BasePrice]++
	}

	result := make([]markets.Pair, 0, len(countMap))

	for price, quantity := range countMap {
		result = append(result, markets.Pair{
			Price:    price,
			Quantity: quantity,
		})
	}

	u.SortPairs(result)

	return result
}
