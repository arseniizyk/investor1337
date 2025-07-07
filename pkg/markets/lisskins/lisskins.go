package lisskins

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

func (ls lisskins) FindByHashName(ctx context.Context, name string) (map[float64]int, error) {
	endpoint := "https://api.lis-skins.com/v1/market/search"
	params := url.Values{
		"game":    []string{"csgo"},
		"names[]": []string{name},
	}

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+ls.token)

	if err != nil {
		ls.l.Error("cant request lis-skins",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	r, err := u.DoJSONRequest[Response](ctx, ls.client, req, ls.l)

	if err != nil {
		ls.l.Warn("response error from lis-skins",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	result := make(map[float64]int, 1)

	for _, o := range r.Data {
		if len(result) == markets.MaxOutputs {
			break
		}
		result[o.Price]++
	}

	return result, nil
}
