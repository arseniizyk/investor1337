package steam

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type steam struct {
	client *http.Client
	items  map[string]int
	l      *zap.Logger
}

type Response struct {
	Success        int    `json:"success"`
	SellOrderCount string `json:"sell_order_count"`
	SellOrderPrice string `json:"sell_order_price"`
	SellOrderTable []struct {
		Price        string `json:"price"`
		PriceWithFee string `json:"price_with_fee"`
		Quantity     string `json:"quantity"`
	} `json:"sell_order_table"`
	BuyOrderCount string `json:"buy_order_count"`
	BuyOrderPrice string `json:"buy_order_price"`
	BuyOrderTable []struct {
		Price    string `json:"price"`
		Quantity string `json:"quantity"`
	} `json:"buy_order_table"`
	HighestBuyOrder string  `json:"highest_buy_order"`
	LowestSellOrder string  `json:"lowest_sell_order"`
	BuyOrderGraph   [][]any `json:"buy_order_graph"`
	SellOrderGraph  [][]any `json:"sell_order_graph"`
	GraphMaxY       int     `json:"graph_max_y"`
	GraphMinX       float64 `json:"graph_min_x"`
	GraphMaxX       float64 `json:"graph_max_x"`
	PricePrefix     string  `json:"price_prefix"`
	PriceSuffix     string  `json:"price_suffix"`
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
	file, err := os.ReadFile("../cs2ids.json")
	if err != nil {
		s.l.Error("Cant load cs2 ids from json", zap.Error(err))
		return err
	}

	data := make(map[string]int)

	if err := json.Unmarshal(file, &data); err != nil {
		s.l.Error("Cant unmarshal cs2 ids", zap.Error(err))
		return err
	}

	s.items = make(map[string]int, len(data))
	for k, v := range data {
		lower := strings.ToLower(k)
		s.items[lower] = v
	}

	return nil
}
