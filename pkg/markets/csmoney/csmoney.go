package csmoney

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type csmoney struct {
	logger *zap.Logger
}

func New(l *zap.Logger) markets.Market {
	return csmoney{l}
}

func (csm csmoney) FindByHashName(name string) (map[float64]int, error) {
	name = url.QueryEscape(name)
	url := fmt.Sprintf("https://cs.money/2.0/market/sell-orders?limit=60&offset=0&name=%s&order=asc&sort=price", name)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		csm.logger.Error("Cant make request", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		csm.logger.Error("cs.money status code bad request:", zap.Int("status_code", resp.StatusCode), zap.String("name", name), zap.String("url", url))
		return nil, err

	case http.StatusOK, http.StatusNotModified:
		var res Response

		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			csm.logger.Error("cant decode cs.money response", zap.Int("status_code", resp.StatusCode), zap.String("name", name))
			return nil, err
		}

		csm.logger.Debug("formatting cs.money response")
		results := format(res)

		csm.logger.Info("csmoney FindByHashName success", zap.String("name", name))
		return results, nil

	default:
		csm.logger.Error("invalid status code", zap.Int("status_code", resp.StatusCode))
		return nil, nil
	}
}

func format(r Response) map[float64]int {
	results := make(map[float64]int)

	price := r.Items[0].Pricing.BasePrice
	count := 1

	for _, item := range r.Items {
		if len(results) == markets.MaxOutputs-1 {
			break
		}
		if item.Pricing.BasePrice != price {
			results[price] = count
			count = 1
			price = item.Pricing.BasePrice
			continue
		}
		count++
	}

	results[price] = count

	return results
}
