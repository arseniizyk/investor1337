package lisskins

import (
	"net/http"
	"time"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type Response struct {
	Data []struct {
		ID             int       `json:"id"`
		Name           string    `json:"name"`
		Price          float64   `json:"price"`
		UnlockAt       any       `json:"unlock_at"`
		ItemClassID    string    `json:"item_class_id"`
		CreatedAt      time.Time `json:"created_at"`
		ItemAssetID    string    `json:"item_asset_id"`
		GameID         int       `json:"game_id"`
		ItemFloat      any       `json:"item_float"`
		NameTag        any       `json:"name_tag"`
		ItemPaintIndex any       `json:"item_paint_index"`
		ItemPaintSeed  any       `json:"item_paint_seed"`
		Stickers       []any     `json:"stickers"`
	} `json:"data"`
	Meta struct {
		PerPage    int    `json:"per_page"`
		NextCursor string `json:"next_cursor"`
	} `json:"meta"`
}

type lisskins struct {
	client *http.Client
	token  string
	l      *zap.Logger
}

func New(client *http.Client, token string, l *zap.Logger) markets.Market {
	return &lisskins{client, token, l}
}
