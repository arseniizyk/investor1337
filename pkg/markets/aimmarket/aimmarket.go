package aimmarket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	m "github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

func (am aimmarket) FindByHashName(ctx context.Context, name string) ([]m.Pair, error) {
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
		return nil, m.ErrRequestFailed
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	r, err := m.DoJSONRequest[Response](ctx, am.client, req, am.l)
	if err != nil {
		if errors.Is(err, m.ErrNoOffers) {
			am.l.Warn("No offers for aimmarket", zap.String("name", name))
			return nil, err
		}

		am.l.Warn("Response error from aim.market",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, m.ErrBadResponse
	}

	result, err := format(&r)
	if err != nil {
		am.l.Warn("cant format output from aimmarket, no offers",
			zap.String("name", name),
		)
		return nil, m.ErrFormatFailed
	}

	return result, nil
}

func (am aimmarket) URL(name string) string {
	return "https://aim.market/ru/buy/csgo/" + url.PathEscape(name)
}

func format(r *Response) ([]m.Pair, error) {
	res := r.Data.BotsInventoryCountAndMinPrice
	if len(res) == 0 {
		return nil, m.ErrBadResponse
	}

	p := r.Data.BotsInventoryCountAndMinPrice[0].Price.SellPrice
	count := r.Data.BotsInventoryCountAndMinPrice[0].Count

	return m.SinglePair(p, count), nil
}
