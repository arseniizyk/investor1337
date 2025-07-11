package steam

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

func (s steam) FindByHashName(ctx context.Context, name string) ([]markets.Pair, error) {
	url := fmt.Sprintf("https://steamcommunity.com/market/itemordershistogram?norender=1&language=english&currency=1&item_nameid=%d", s.items[strings.ToLower(name)])

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		s.l.Error("Cant make request to steam",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	r, err := u.DoJSONRequest[Response](ctx, s.client, req, s.l)
	if err != nil {
		s.l.Warn("Response error from steam",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	results, err := format(&r)
	if err != nil {
		s.l.Error("Cant format response from steam",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	return results, nil
}

func format(r *Response) ([]markets.Pair, error) {
	results := make([]markets.Pair, 0, markets.MaxOutputs)
	re := regexp.MustCompile(`\d{1,3}(?:,\d{3})*(?:\.\d+)?|\d+\.\d+|\d+`)

	for i, orders := range r.SellOrderTable {
		if i == markets.MaxOutputs {
			break
		}

		matched := re.FindString(orders.Price)
		priceStr := strings.ReplaceAll(matched, ",", "")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return nil, err
		}

		quantity, err := strconv.Atoi(strings.ReplaceAll(orders.Quantity, ",", ""))
		if err != nil {
			return nil, err
		}

		results = append(results, markets.Pair{
			Price:    price,
			Quantity: quantity,
		})
	}

	u.SortPairs(results)

	return results, nil
}
