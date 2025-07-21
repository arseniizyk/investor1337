package steam

import (
	"embed"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

//go:embed cs2ids.json
var cs2ids embed.FS

type steam struct {
	client *http.Client
	items  map[string]int
	l      *zap.Logger
}

type Response struct {
	Success        int `json:"success"`
	SellOrderTable []struct {
		Price        string `json:"price"`
		PriceWithFee string `json:"price_with_fee"`
		Quantity     string `json:"quantity"`
	} `json:"sell_order_table"`
}

func New(c *http.Client, l *zap.Logger) (markets.Market, error) {
	s := steam{
		client: c,
		l:      l,
	}

	if err := s.loadItems(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *steam) loadItems() error {
	b, err := cs2ids.ReadFile("cs2ids.json")
	if err != nil {
		s.l.Error("Cant load cs2 ids from json", zap.Error(err))
		return err
	}

	data := make(map[string]int)

	if err := json.Unmarshal(b, &data); err != nil {
		s.l.Error("Cant unmarshal cs2 ids", zap.Error(err))
		return err
	}

	s.items = make(map[string]int, len(data))
	for k, v := range data {
		s.items[strings.ToLower(k)] = v
	}

	return nil
}
