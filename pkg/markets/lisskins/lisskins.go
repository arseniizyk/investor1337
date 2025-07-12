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

func (ls lisskins) buildRequest(cursor, name string) (*http.Request, error) {
	endpoint := "https://api.lis-skins.com/v1/market/search"
	params := url.Values{
		"game":    []string{"csgo"},
		"names[]": []string{name},
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
		return nil, u.ErrRequest
	}

	return req, nil
}

func (ls lisskins) FindByHashName(ctx context.Context, name string) ([]markets.Pair, error) {
	countMap, err := u.FetchWithCursor(ctx, ls.client, ls.l, name, "LIS-SKINS", maxPages, countInMap, ls.buildRequest)
	if err != nil {
		ls.l.Warn("LIS-SKINS error in FetchWithCursor",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, err
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
