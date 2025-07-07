package csfloat

import (
	"net/http"
	"time"

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
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Type      string    `json:"type"`
		Price     int       `json:"price"`
		State     string    `json:"state"`
		Seller    struct {
			Away                bool   `json:"away"`
			Flags               int    `json:"flags"`
			HasValidSteamAPIKey bool   `json:"has_valid_steam_api_key"`
			ObfuscatedID        string `json:"obfuscated_id"`
			Online              bool   `json:"online"`
			StallPublic         bool   `json:"stall_public"`
			Statistics          struct {
				MedianTradeTime     int `json:"median_trade_time"`
				TotalAvoidedTrades  int `json:"total_avoided_trades"`
				TotalFailedTrades   int `json:"total_failed_trades"`
				TotalTrades         int `json:"total_trades"`
				TotalVerifiedTrades int `json:"total_verified_trades"`
			} `json:"statistics"`
		} `json:"seller"`
		Reference struct {
			BasePrice      int       `json:"base_price"`
			PredictedPrice int       `json:"predicted_price"`
			Quantity       int       `json:"quantity"`
			LastUpdated    time.Time `json:"last_updated"`
		} `json:"reference"`
		Item struct {
			AssetID        string `json:"asset_id"`
			DefIndex       int    `json:"def_index"`
			IconURL        string `json:"icon_url"`
			Rarity         int    `json:"rarity"`
			MarketHashName string `json:"market_hash_name"`
			Tradable       int    `json:"tradable"`
			IsCommodity    bool   `json:"is_commodity"`
			Type           string `json:"type"`
			RarityName     string `json:"rarity_name"`
			TypeName       string `json:"type_name"`
			ItemName       string `json:"item_name"`
		} `json:"item"`
		IsSeller       bool `json:"is_seller"`
		IsWatchlisted  bool `json:"is_watchlisted"`
		Watchers       int  `json:"watchers"`
		AuctionDetails struct {
			ReservePrice int `json:"reserve_price"`
			TopBid       struct {
				ID                string    `json:"id"`
				CreatedAt         time.Time `json:"created_at"`
				Price             int       `json:"price"`
				ContractID        string    `json:"contract_id"`
				State             string    `json:"state"`
				ObfuscatedBuyerID string    `json:"obfuscated_buyer_id"`
			} `json:"top_bid"`
			ExpiresAt  time.Time `json:"expires_at"`
			MinNextBid int       `json:"min_next_bid"`
		} `json:"auction_details"`
	} `json:"data"`
	Cursor string `json:"cursor"`
}

func New(c *http.Client, cookie string, l *zap.Logger) markets.Market {
	return &csfloat{c, cookie, l}
}
