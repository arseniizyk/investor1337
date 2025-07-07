package steam

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	u "github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

func (s steam) FindByHashName(ctx context.Context, name string) (map[float64]int, error) {
	url := fmt.Sprintf("https://steamcommunity.com/market/itemordershistogram?norender=1&language=english&currency=1&item_nameid=%d", s.data[strings.ToLower(name)])

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		s.l.Error("cant make request to steam",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	r, err := u.DoJSONRequest[Response](ctx, s.client, req, s.l)

	if err != nil {
		s.l.Warn("response error from steam",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	results, err := format(&r)
	if err != nil {
		s.l.Error("cant format response from steam",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	return results, nil
}

func format(r *Response) (map[float64]int, error) {
	results := make(map[float64]int)

	for i, orders := range r.SellOrderTable {
		if i == 3 {
			break
		}

		re := regexp.MustCompile(`\d{1,3}(?:,\d{3})*(?:\.\d+)?|\d+\.\d+|\d+`)
		matched := re.FindString(orders.Price)
		priceString := strings.ReplaceAll(matched, ",", "")

		price, err := strconv.ParseFloat(priceString, 64)
		if err != nil {
			return nil, err
		}

		quantity, err := strconv.Atoi(strings.ReplaceAll(orders.Quantity, ",", ""))
		if err != nil {
			return nil, err
		}

		results[price] = quantity
	}

	return results, nil
}
