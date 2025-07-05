package csgomarket

import (
	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

const commissionMultiplier = 1.038

type csgoMarket struct {
	token string
	l     *zap.Logger
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

func New(token string, l *zap.Logger) markets.Market {
	return csgoMarket{token, l}
}
