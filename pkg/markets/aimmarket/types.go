package aimmarket

import (
	"embed"
	"net/http"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

//go:embed query.graphql
var graphqlQuery embed.FS

type aimmarket struct {
	client *http.Client
	l      *zap.Logger
	query  []byte
}

type Response struct {
	Data struct {
		BotsInventoryCountAndMinPrice []struct {
			MarketHashName string `json:"marketHashName"`
			Count          int    `json:"count"`
			Price          struct {
				SellPrice float64 `json:"sellPrice"`
				Currency  string  `json:"currency"`
				Typename  string  `json:"__typename"`
			} `json:"price"`
			Typename string `json:"__typename"`
		} `json:"bots_inventory_count_and_min_price"`
	} `json:"data"`
}

func New(c *http.Client, l *zap.Logger) (markets.Market, error) {
	am := aimmarket{
		client: c,
		l:      l,
	}

	if err := am.loadGraphQlQuery(); err != nil {
		return nil, err
	}

	return am, nil
}

func (am *aimmarket) loadGraphQlQuery() error {
	data, err := graphqlQuery.ReadFile("query.graphql")
	if err != nil {
		am.l.Error("Cant load GraphQL query file", zap.Error(err))
		return err
	}

	am.query = data
	return nil
}
