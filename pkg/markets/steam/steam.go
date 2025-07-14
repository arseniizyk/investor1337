package steam

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	m "github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

func (s steam) FindByHashName(ctx context.Context, name string) ([]m.Pair, error) {
	url := fmt.Sprintf("https://steamcommunity.com/market/itemordershistogram?norender=1&language=english&currency=1&item_nameid=%d", s.items[strings.ToLower(name)])

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		s.l.Error("Cant make request to steam",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, m.ErrRequestFailed
	}

	r, err := m.DoJSONRequest[Response](ctx, s.client, req, s.l)
	if err != nil {
		s.l.Warn("Response error from steam",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, m.ErrBadResponse
	}

	results, err := format(&r)
	if err != nil {
		s.l.Error("Cant format response from steam",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, m.ErrFormatFailed
	}

	return results, nil
}

func (s steam) URL(name string) string {
	return "https://steamcommunity.com/market/listings/730/" + url.PathEscape(name)
}

func format(r *Response) ([]m.Pair, error) {
	results := make([]m.Pair, 0, m.MaxOutputs)
	re := regexp.MustCompile(`\d{1,3}(?:,\d{3})*(?:\.\d+)?|\d+\.\d+|\d+`)

	for i, orders := range r.SellOrderTable {
		if i == m.MaxOutputs {
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

		results = append(results, m.Pair{
			Price:    price,
			Quantity: quantity,
		})
	}

	m.SortPairs(results)

	return results, nil
}
