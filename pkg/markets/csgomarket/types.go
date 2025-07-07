package csgomarket

import (
	"net/http"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type csgoMarket struct {
	client *http.Client
	token  string
	l      *zap.Logger
}

type Response struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Currency string `json:"currency"`
	Data     []struct {
		MarketHashName string `json:"market_hash_name"`
		Price          int    `json:"price"`
		Class          int    `json:"class"`
		Instance       int    `json:"instance"`
		Count          int    `json:"count"`
	} `json:"data"`
}

func New(c *http.Client, token string, l *zap.Logger) markets.Market {
	return csgoMarket{c, token, l}
}
