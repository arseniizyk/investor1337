package aimmarket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

func (am aimmarket) FindByHashName(name string) (map[float64]int, error) {
	payload := map[string]any{
		"operationName": "ApiBotsInventoryCountAndMinPrice",
		"query":         string(am.query),
		"variables": map[string]any{
			"currency": "USD",
			"where": map[string]any{
				"marketHashName": map[string]string{
					"_text": fmt.Sprintf("\"%s\"", name),
				},
			},
		},
	}

	b, err := json.Marshal(payload)
	if err != nil {
		am.l.Error("cant marshal payload in aimmarket",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://aim.market/v1/api/graphql", bytes.NewBuffer(b))
	if err != nil {
		am.l.Error("cant request aimmarket",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		am.l.Error("bad response from aimmarket",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}
	defer utils.Dclose(resp.Body, am.l)

	switch resp.StatusCode {
	case http.StatusOK:
		var r Response

		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			am.l.Error("cant decode response from aimmarket",
				zap.String("name", name),
				zap.Int("status_code", resp.StatusCode),
				zap.Error(err),
			)
			return nil, err
		}

		res := r.Data.BotsInventoryCountAndMinPrice
		if len(res) == 0 {
			am.l.Warn("No offers for aimmarket", zap.String("name", name))
			return nil, errors.New("no offers")
		}

		p := r.Data.BotsInventoryCountAndMinPrice[0].Price.SellPrice
		count := r.Data.BotsInventoryCountAndMinPrice[0].Count

		result := map[float64]int{p: count}

		return result, nil
	default:
		am.l.Warn("unknown status_code",
			zap.Int("status_code", resp.StatusCode),
			zap.String("name", name))
		return nil, errors.New("bad response code")
	}
}
