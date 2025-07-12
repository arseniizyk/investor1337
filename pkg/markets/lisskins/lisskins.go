package lisskins

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/utils"
	"go.uber.org/zap"
)

const maxPages = 5 // 200 items per page

func (ls lisskins) FindByHashName(ctx context.Context, name string) ([]markets.Pair, error) {
	endpoint := "https://api.lis-skins.com/v1/market/search"
	params := url.Values{
		"game":    []string{"csgo"},
		"names[]": []string{name},
	}

	countMap := make(map[float64]int, markets.MaxOutputs)
	cursor := ""

	for range maxPages {
		if len(countMap) == markets.MaxOutputs {
			break
		}

		if cursor != "" {
			params.Set("cursor", cursor)
		}

		url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

		req, err := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+ls.token)

		if err != nil {
			ls.l.Error("Cant make request to lis-skins",
				zap.String("name", name),
				zap.Error(err),
			)
			return nil, err
		}

		r, err := u.DoJSONRequest[Response](ctx, ls.client, req, ls.l)

		if err != nil {
			ls.l.Warn("Response error from lis-skins",
				zap.String("name", name),
				zap.Error(err),
			)
			return nil, err
		}

		if len(r.Data) == 0 {
			break
		}

		countInMap(countMap, &r)

		ls.l.Debug("LIS-SKINS Fetched page",
			zap.String("name", name),
			zap.Int("items", len(r.Data)),
			zap.String("cursor", cursor),
		)

		cursor = r.Meta.NextCursor

		if r.Meta.NextCursor == "" {
			break
		}
	}

	return u.PairsFromMap(countMap), nil
}

func (ls lisskins) URL(name string) string {
	var last rune
	formatted := strings.Map(func(r rune) rune {
		switch r {
		case '(', ')', '|', ':', ',', '.', '’', '"', '\'':
			return -1
		case ' ':
			if last == '-' {
				return -1 // to avoid double -
			}
			last = '-'
			return '-'
		default:
			lower := unicode.ToLower(r)
			last = lower
			return lower
		}
	}, name)

	formatted = strings.Trim(formatted, "-")

	return "https://lis-skins.com/ru/market/csgo/" + formatted
}

func countInMap(m map[float64]int, r *Response) {
	for _, o := range r.Data {
		if len(m) == markets.MaxOutputs {
			break
		}
		m[o.Price]++
	}
}
