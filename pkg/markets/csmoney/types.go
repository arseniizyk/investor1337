package csmoney

import (
	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

const commissionMult = 1.042

type csmoney struct {
	l *zap.Logger
}

type Response struct {
	Items []struct {
		ID     int `json:"id"`
		AppID  int `json:"appId"`
		Seller struct {
			BotID    any `json:"botId"`
			Delivery struct {
				Speed       string  `json:"speed"`
				MedianTime  float64 `json:"medianTime"`
				SuccessRate float64 `json:"successRate"`
			} `json:"delivery"`
		} `json:"seller"`
		Asset struct {
			ID    int64 `json:"id"`
			Names struct {
				Short      string `json:"short"`
				Full       string `json:"full"`
				Identifier int    `json:"identifier"`
			} `json:"names"`
			Images struct {
				Steam string `json:"steam"`
			} `json:"images"`
			IsSouvenir bool   `json:"isSouvenir"`
			IsStatTrak bool   `json:"isStatTrak"`
			Rarity     string `json:"rarity"`
			Pattern    any    `json:"pattern"`
			Type       int    `json:"type"`
			Collection struct {
				Name  string `json:"name"`
				Image any    `json:"image"`
			} `json:"collection"`
			Float   any `json:"float"`
			Inspect any `json:"inspect"`
		} `json:"asset"`
		Stickers  any `json:"stickers"`
		Keychains any `json:"keychains"`
		Pricing   struct {
			Default             float64 `json:"default"`
			PriceBeforeDiscount float64 `json:"priceBeforeDiscount"`
			Discount            float64 `json:"discount"`
			Computed            float64 `json:"computed"`
			BasePrice           float64 `json:"basePrice"`
			PriceCoefficient    float64 `json:"priceCoefficient"`
		} `json:"pricing"`
		Links struct {
			ThreeD      any    `json:"3d"`
			InspectLink string `json:"inspectLink"`
		} `json:"links"`
		IsPartial     bool `json:"isPartial"`
		IsMySellOrder bool `json:"isMySellOrder"`
	} `json:"items"`
}

func New(l *zap.Logger) markets.Market {
	return csmoney{l}
}
