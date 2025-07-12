package lisskins

import (
	"net/http"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type Response struct {
	Data []struct {
		Price float64 `json:"price"`
	} `json:"data"`
	Meta struct {
		NextCursor string `json:"next_cursor"`
	} `json:"meta"`
}

func (r Response) LenData() int {
	return len(r.Data)
}

func (r Response) Cursor() string {
	return r.Meta.NextCursor
}

type lisskins struct {
	client *http.Client
	token  string
	l      *zap.Logger
}

func New(c *http.Client, token string, l *zap.Logger) markets.Market {
	return &lisskins{c, token, l}
}
