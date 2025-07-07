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

func (csm csmoney) FindByHashName(ctx context.Context, name string) (map[float64]int, error) {
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

func format(r *Response) map[float64]int {
	results := make(map[float64]int)

	price := r.Items[0].Pricing.BasePrice
	count := 1

	for _, item := range r.Items {
		if len(results) == markets.MaxOutputs {
			break
		}
		if item.Pricing.BasePrice != price {
			results[price] = count
			count = 1
			price = item.Pricing.BasePrice
			continue
		}
		count++
	}

	results[price] = count

	return results
}
