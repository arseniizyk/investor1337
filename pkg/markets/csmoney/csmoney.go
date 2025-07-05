package csmoney

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

func (csm csmoney) FindByHashName(name string) (map[float64]int, error) {
	url := fmt.Sprintf("https://cs.money/2.0/market/sell-orders?limit=60&offset=0&name=%s&order=asc&sort=price", url.QueryEscape(name))

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	if err != nil {
		csm.l.Error("cant make request to cs.money",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		csm.l.Error("cant request cs.money",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}
	defer utils.Dclose(resp.Body, csm.l)

	switch resp.StatusCode {
	case http.StatusOK, http.StatusNotModified:
		var r Response

		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			csm.l.Error("cant decode response from cs.money",
				zap.String("name", name),
				zap.Error(err))
			return nil, err
		}

		result := format(&r)

		return result, nil

	case http.StatusBadRequest:
		csm.l.Warn("cs.money status code bad request", zap.String("name", name))
		return nil, errors.New("bad request")

	default:
		csm.l.Warn("unknown status code from cs.money",
			zap.Int("status_code", resp.StatusCode),
			zap.String("name", name))
		return nil, errors.New("unknown status code")
	}
}

func format(r *Response) map[float64]int {
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
