package csfloat

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

// because csfloat returns price in int, not float64
const priceDivider = 100.0

func (c csfloat) FindByHashName(ctx context.Context, name string) (map[float64]int, error) {
	endpoint := "https://csfloat.com/api/v1/listings"
	params := url.Values{
		"limit":            []string{"0"},
		"market_hash_name": []string{name},
		"sort_by":          []string{"lowest_price"},
	}

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.l.Error("Cant make request to csfloat",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	req.Header.Add("Cookie", c.cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	r, err := u.DoJSONRequest[Response](ctx, c.client, req, c.l)
	if err != nil {
		c.l.Warn("Response error from csfloat",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	result := make(map[float64]int, 1)

	for _, seller := range r.Data {
		if len(result) == markets.MaxOutputs {
			break
		}

		p := float64(seller.Price) / priceDivider
		result[p]++
	}

	return result, nil
}
