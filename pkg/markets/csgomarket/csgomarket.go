package csgomarket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/utils"
	"go.uber.org/zap"
)

// because csgo market returns int, not float64 so we need to divide it by 1000
const priceDivider = 1000.0

func (cm csgoMarket) FindByHashName(ctx context.Context, name string) ([]markets.Pair, error) {
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
		return nil, err
	}

	r, err := u.DoJSONRequest[Response](ctx, cm.client, req, cm.l)
	if err != nil {
		cm.l.Warn("Response error from csgo market",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	if !r.Success || len(r.Data) == 0 {
		cm.l.Warn("Csgo market bad request", zap.String("name", name))
		return nil, errors.New("bad request")
	}

	return format(&r), nil
}

func format(r *Response) []markets.Pair {
	result := make([]markets.Pair, 0, markets.MaxOutputs)

	for _, o := range r.Data {
		if len(result) == markets.MaxOutputs {
			break
		}
		p := float64(o.Price) / priceDivider
		result = append(result, markets.Pair{
			Price:    p,
			Quantity: o.Count,
		})
	}

	u.SortPairs(result)

	return result
}
