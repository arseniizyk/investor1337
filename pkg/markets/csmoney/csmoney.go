package csmoney

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

const commissionCoeff = 1.042

type csmoney struct {
	l *zap.Logger
}

func New(l *zap.Logger) markets.Market {
	return csmoney{l}
}

func (csm csmoney) FindByHashName(name string) (map[float64]int, error) {
	name = url.QueryEscape(name)
	url := fmt.Sprintf("https://cs.money/2.0/market/sell-orders?limit=60&offset=0&name=%s&order=asc&sort=price", name)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	if err != nil {
		csm.l.Error("Cant make request for csmoney", zap.Error(err))
		return nil, err
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		csm.l.Error("Cant make request", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	defer utils.Dclose(resp.Body, csm.l)

	switch resp.StatusCode {
	case http.StatusBadRequest:
		csm.l.Error("cs.money status code bad request:", zap.Int("status_code", resp.StatusCode), zap.String("name", name), zap.String("url", url))
		return nil, err

	case http.StatusOK, http.StatusNotModified:
		var res Response

		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			csm.l.Error("cant decode cs.money response", zap.Int("status_code", resp.StatusCode), zap.String("name", name))
			return nil, err
		}

		csm.l.Debug("formatting cs.money response")
		results := format(res)

		csm.l.Debug("csmoney FindByHashName success", zap.String("name", name))
		return results, nil

	default:
		csm.l.Error("invalid status code", zap.Int("status_code", resp.StatusCode))
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
			results[price*commissionCoeff] = count
			count = 1
			price = item.Pricing.BasePrice
			continue
		}
		count++
	}

	results[price] = count

	return results
}
