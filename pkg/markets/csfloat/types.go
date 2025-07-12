package csfloat

import (
	"net/http"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type csfloat struct {
	client *http.Client
	cookie string
	l      *zap.Logger
}

type Response struct {
	Data []struct {
		Price int `json:"price"`
	} `json:"data"`
	Cursor string `json:"cursor"`
}

func New(c *http.Client, cookie string, l *zap.Logger) markets.Market {
	return &csfloat{c, cookie, l}
}
