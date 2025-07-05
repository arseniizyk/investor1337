package lisskins

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

func (ls lisskins) FindByHashName(name string) (map[float64]int, error) {
	endpoint := "https://api.lis-skins.com/v1/market/search"
	params := url.Values{}
	params.Set("game", "csgo")
	params.Add("names[]", name)

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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ls.l.Error("response error from lis-skins",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}
	defer utils.Dclose(resp.Body, ls.l)

	switch resp.StatusCode {
	case http.StatusOK:
		var r Response

		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			ls.l.Error("cant decode response from lis-skins",
				zap.String("name", name),
				zap.Int("status_code", resp.StatusCode),
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

	case http.StatusUnauthorized:
		ls.l.Error("lis-skins bad token, generate a new one", zap.String("token", ls.token))
		return nil, errors.New("bad token")

	case http.StatusUnprocessableEntity:
		ls.l.Warn("lis-skins bad request", zap.String("name", name))
		return nil, errors.New("bad request")

	default:
		ls.l.Warn("unknown status from lis-skins",
			zap.Int("status_code", resp.StatusCode),
			zap.String("name", name))

		return nil, errors.New("unknown status code")
	}
}
