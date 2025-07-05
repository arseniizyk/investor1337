package csgomarket

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

func (cm csgoMarket) FindByHashName(name string) (map[float64]int, error) {
	endpoint := fmt.Sprintf("https://market.csgo.com/api/v2/search-item-by-hash-name?key=%s&hash_name=%s", cm.token, url.QueryEscape(name))

	resp, err := http.Get(endpoint)
	if err != nil {
		cm.l.Error("cant request csgo market",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}
	defer utils.Dclose(resp.Body, cm.l)

	switch resp.StatusCode {
	case http.StatusOK:
		var r Response

		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			cm.l.Error("cant decode response from csgo market",
				zap.String("name", name),
				zap.Int("status_code", resp.StatusCode),
				zap.Error(err))
			return nil, err
		}

		if !r.Success || len(r.Data) == 0 {
			cm.l.Warn("csgo market bad request",
				zap.String("name", name),
				zap.Int("status_code", resp.StatusCode))
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

	case http.StatusBadRequest:
		cm.l.Warn("csgo market bad request", zap.String("name", name))
		return nil, errors.New("bad request")

	default:
		cm.l.Warn("unknown status from csgo market",
			zap.Int("status_code", resp.StatusCode),
			zap.String("name", name))

		return nil, errors.New("unknown status code")
	}
}
