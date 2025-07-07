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

func (b buff163) FindByHashName(ctx context.Context, name string) (map[float64]int, error) {
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
			zap.Error(err))
		return nil, err
	}

	r, err := u.DoJSONRequest[Response](ctx, b.client, req, b.l)
	if err != nil {
		b.l.Warn("Response error from buff163",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	result := make(map[float64]int, 1)

	for _, i := range r.Data.Items {
		if len(result) == markets.MaxOutputs {
			break
		}

		price, err := strconv.ParseFloat(i.Price, 64)
		if err != nil {
			b.l.Error("Cant parse to float64",
				zap.String("price", i.Price),
				zap.Error(err))
			return nil, err
		}

		p := math.Round(price*priceMultiplier*100) / 100
		result[p]++
	}

	return result, nil
}
