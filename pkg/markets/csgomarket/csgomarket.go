package csgomarket

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	m "github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

// because csgo market returns int, not float64 so we need to divide it by 1000
const priceDivider = 1000.0

func (cm csgoMarket) FindByHashName(ctx context.Context, name string) ([]m.Pair, error) {
	endpoint := "https://market.csgo.com/api/v2/search-item-by-hash-name"
	params := url.Values{
		"key":       []string{cm.token},
		"hash_name": []string{name},
	}

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		cm.l.Error("Cant make request to csgo market",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, m.ErrRequestFailed
	}

	r, err := m.DoJSONRequest[Response](ctx, cm.client, req, cm.l)
	if err != nil {
		cm.l.Warn("Response error from csgo market",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, m.ErrBadResponse
	}

	if !r.Success || len(r.Data) == 0 {
		cm.l.Warn("Csgo market bad request", zap.String("name", name))
		return nil, m.ErrNoOffers
	}

	return format(&r), nil
}

func (cm csgoMarket) URL(name string) string {
	return "https://market.csgo.com/en/" + url.PathEscape(name)
}

func format(r *Response) []m.Pair {
	result := make([]m.Pair, 0, m.MaxOutputs)

	for _, o := range r.Data {
		if len(result) == m.MaxOutputs {
			break
		}
		p := float64(o.Price) / priceDivider
		result = append(result, m.Pair{
			Price:    p,
			Quantity: o.Count,
		})
	}

	m.SortPairs(result)

	return result
}
