package buff163

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

const priceMultiplier = 0.14 // convert CNY to USD

func (b buff163) FindByHashName(ctx context.Context, name string) ([]markets.Pair, error) {
	endpoint := "https://buff.163.com/api/market/goods/sell_order"
	goodsId := strconv.Itoa(b.items[strings.ToLower(name)])
	params := url.Values{
		"game":     []string{"csgo"},
		"goods_id": []string{goodsId},
	}

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		b.l.Error("Cant make request to buff163",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	r, err := u.DoJSONRequest[Response](ctx, b.client, req, b.l)
	if err != nil {
		b.l.Warn("Response error from buff163",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	result, err := format(&r)
	if err != nil {
		b.l.Error("cant format output from buff163",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}

func format(r *Response) ([]markets.Pair, error) {
	countMap := make(map[float64]int, 1)

	for _, i := range r.Data.Items {
		if len(countMap) == markets.MaxOutputs {
			break
		}

		price, err := strconv.ParseFloat(i.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("cant parse to float64: price: %s, err: %v", i.Price, err)
		}

		p := math.Round(price*priceMultiplier*100) / 100
		countMap[p]++
	}

	return u.PairsFromMap(countMap), nil
}
