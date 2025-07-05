package steam

import (
	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type steam struct {
	data map[string]int
	l    *zap.Logger
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
	HighestBuyOrder string          `json:"highest_buy_order"`
	LowestSellOrder string          `json:"lowest_sell_order"`
	BuyOrderGraph   [][]interface{} `json:"buy_order_graph"`
	SellOrderGraph  [][]interface{} `json:"sell_order_graph"`
	GraphMaxY       int             `json:"graph_max_y"`
	GraphMinX       float64         `json:"graph_min_x"`
	GraphMaxX       float64         `json:"graph_max_x"`
	PricePrefix     string          `json:"price_prefix"`
	PriceSuffix     string          `json:"price_suffix"`
}

func New(logger *zap.Logger) (markets.Market, error) {
	s := steam{
		l: logger,
	}

	if err := s.loadNameIds(); err != nil {
		return nil, err
	}

	return s, nil
}
