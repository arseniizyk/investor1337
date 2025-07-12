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
	maxPages     = 5 // 50 items per page
	priceDivider = 100.0
)

func (c csfloat) buildRequest(cursor, name string) (*http.Request, error) {
	endpoint := "https://csfloat.com/api/v1/listings"
	params := url.Values{
		"limit":            []string{"0"},
		"market_hash_name": []string{name},
		"sort_by":          []string{"lowest_price"},
	}

	if cursor != "" {
		params.Set("cursor", cursor)
	}

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.l.Error("Cant make request to CSFloat",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, u.ErrRequest
	}

	req.Header.Add("Cookie", c.cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	return req, nil
}

func (c csfloat) FindByHashName(ctx context.Context, name string) ([]markets.Pair, error) {
	countMap, err := u.FetchWithCursor(ctx, c.client, c.l, name, "CSFloat", maxPages, countInMap, c.buildRequest)

	if err != nil {
		c.l.Warn("CSFloat error in FetchWithCursor",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
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
