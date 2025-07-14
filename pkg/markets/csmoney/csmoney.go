package csmoney

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	m "github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

func (csm csmoney) FindByHashName(ctx context.Context, name string) ([]m.Pair, error) {
	endpoint := "https://cs.money/2.0/market/sell-orders"
	params := url.Values{
		"limit":  []string{"60"},
		"offset": []string{"0"},
		"name":   []string{name},
		"order":  []string{"asc"},
		"sort":   []string{"price"},
	}

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	if err != nil {
		csm.l.Error("Cant make request to cs.money",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, m.ErrRequestFailed
	}

	r, err := m.DoJSONRequest[Response](ctx, csm.client, req, csm.l)

	if err != nil {
		csm.l.Warn("Response error from cs.money",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, m.ErrBadResponse
	}

	return format(&r), nil
}

func (csm csmoney) URL(name string) string {
	var last rune
	formatted := strings.Map(func(r rune) rune {
		switch r {
		case '(', ')', '|', '&', ':', ',', '.', '’', '"', '\'', '★':
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

	return "https://cs.money/ru/csgo/" + formatted

	// return "https://cs.money/ru/market/buy/" for simplicity, because the code above doesn't work for all items(like doppler knives)
}

func format(r *Response) []m.Pair {
	countMap := make(map[float64]int, m.MaxOutputs)

	for _, item := range r.Items {
		if len(countMap) == m.MaxOutputs {
			break
		}
		countMap[item.Pricing.BasePrice]++
	}

	return m.PairsFromMap(countMap)
}
