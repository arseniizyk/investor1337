package csgomarket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

func (cm csgoMarket) FindByHashName(ctx context.Context, name string) (map[float64]int, error) {
	endpoint := "https://market.csgo.com/api/v2/search-item-by-hash-name"
	params := url.Values{}
	params.Set("key", cm.token)
	params.Set("hash_name", name)

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		cm.l.Error("cant request csgo market",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	r, err := utils.DoJSONRequest[Response](ctx, cm.client, req, cm.l)
	if err != nil {
		cm.l.Warn("response error from csgo market",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	if !r.Success || len(r.Data) == 0 {
		cm.l.Warn("csgo market bad request", zap.String("name", name))
		return nil, errors.New("bad request")
	}

	result := make(map[float64]int, 1)

	for _, o := range r.Data {
		if len(result) == markets.MaxOutputs {
			break
		}
		p := float64(o.Price) / 1000
		result[p] = o.Count
	}

	return result, nil
}
