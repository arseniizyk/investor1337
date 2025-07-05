package csgomarket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

type csgoMarket struct {
	token string
	l     *zap.Logger
}

type Response struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Currency string `json:"currency"`
	Data     []struct {
		MarketHashName string `json:"market_hash_name"`
		Price          int    `json:"price"`
		Class          int    `json:"class"`
		Instance       int    `json:"instance"`
		Count          int    `json:"count"`
	} `json:"data"`
}

func New(token string, l *zap.Logger) markets.Market {
	return csgoMarket{token, l}
}

func (cm csgoMarket) FindByHashName(hashName string) (map[float64]int, error) {
	endpoint := fmt.Sprintf("https://market.csgo.com/api/v2/search-item-by-hash-name?key=%s&hash_name=%s", cm.token, hashName)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("cant get item by hashName: %w", err)
	}
	defer utils.Dclose(resp.Body, cm.l)

	var r Response

	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("cant decode json: statusCode: %d; %w", resp.StatusCode, err)
	}

	if !r.Success {
		return nil, fmt.Errorf("something wrong with response or request: %s", r.Error)
	}

	offers := r.Data

	result := make(map[float64]int, 1)

	for _, o := range offers {
		if len(result) == markets.MaxOutputs {
			break
		}
		p := float64(o.Price) / 1000
		result[p] = o.Count
	}

	return result, nil
}
