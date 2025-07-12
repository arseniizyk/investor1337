package csmoney

import (
	"net/http"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type csmoney struct {
	client *http.Client
	l      *zap.Logger
}

type Response struct {
	Items []struct {
		Pricing struct {
			BasePrice float64 `json:"basePrice"`
		} `json:"pricing"`
	} `json:"items"`
}

func New(c *http.Client, l *zap.Logger) markets.Market {
	return csmoney{c, l}
}
