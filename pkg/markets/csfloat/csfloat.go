package csfloat

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/utils"
	"go.uber.org/zap"
)

// because csfloat returns price in int, not float64
const (
	maxPages     = 7
	priceDivider = 100.0
)

func (c csfloat) FindByHashName(ctx context.Context, name string) ([]markets.Pair, error) {
	endpoint := "https://csfloat.com/api/v1/listings"
	params := url.Values{
		"limit":            []string{"0"},
		"market_hash_name": []string{name},
		"sort_by":          []string{"lowest_price"},
	}

	countMap := make(map[float64]int)
	cursor := ""

	for range maxPages {
		if cursor != "" {
			params.Set("cursor", cursor)
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

		if r.Cursor == "" || len(r.Data) == 0 || len(countMap) == markets.MaxOutputs {
			break
		}

		countInMap(countMap, &r)
		cursor = r.Cursor

		c.l.Debug("CSFloat Fetched page",
			zap.String("name", name),
			zap.Int("items", len(r.Data)),
			zap.String("cursor", r.Cursor),
		)
	}

	return u.PairsFromMap(countMap), nil
}

func (c csfloat) URL(name string) string {
	return "https://csfloat.com/search?market_hash_name=" + url.PathEscape(name)
}

func countInMap(m map[float64]int, r *Response) {
	for _, seller := range r.Data {
		if len(m) == markets.MaxOutputs {
			break
		}

		p := float64(seller.Price) / priceDivider
		m[p]++
	}
}
