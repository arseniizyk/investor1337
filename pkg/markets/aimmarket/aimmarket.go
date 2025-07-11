package aimmarket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/utils"
	"go.uber.org/zap"
)

var ErrNoOffers = errors.New("no offers")

func (am aimmarket) FindByHashName(ctx context.Context, name string) ([]markets.Pair, error) {
	payload := am.preparePayload(name)

	b, err := json.Marshal(payload)
	if err != nil {
		am.l.Error("Cant marshal payload in aimmarket",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://aim.market/v1/api/graphql", bytes.NewBuffer(b))
	if err != nil {
		am.l.Error("Cant make request to aimmarket",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	r, err := u.DoJSONRequest[Response](ctx, am.client, req, am.l)
	if err != nil {
		am.l.Warn("Response error from aim.market",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
	}

	result, err := format(&r)
	if err != nil {
		am.l.Warn("no offers for aimmarket",
			zap.String("name", name),
		)
		return nil, err
	}

	return result, nil
}

func (am aimmarket) URL(name string) string {
	return "https://aim.market/ru/buy/csgo/" + url.PathEscape(name)
}

func format(r *Response) ([]markets.Pair, error) {
	res := r.Data.BotsInventoryCountAndMinPrice
	if len(res) == 0 {
		return nil, ErrNoOffers
	}

	p := r.Data.BotsInventoryCountAndMinPrice[0].Price.SellPrice
	count := r.Data.BotsInventoryCountAndMinPrice[0].Count

	return u.SinglePair(p, count), nil
}
