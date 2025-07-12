package buff163

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type buff163 struct {
	client *http.Client
	items  map[string]int
	l      *zap.Logger
}

type Response struct {
	Code string `json:"code"`
	Data struct {
		Items []struct {
			Price string `json:"price"`
		} `json:"items"`
	} `json:"data"`
}

func New(c *http.Client, l *zap.Logger) (markets.Market, error) {
	buff := buff163{
		client: c,
		l:      l,
	}

	if err := buff.loadItems(); err != nil {
		return nil, err
	}

	return buff, nil
}

func (buff *buff163) loadItems() error {
	b, err := os.ReadFile("../buff163ids.json")
	if err != nil {
		buff.l.Error("Cant load buff163 ids from json", zap.Error(err))
		return err
	}

	data := make(map[string]int)

	if err := json.Unmarshal(b, &data); err != nil {
		buff.l.Error("Cant unmarshal buff163 ids", zap.Error(err))
		return err
	}

	buff.items = make(map[string]int, len(data))
	for k, v := range data {
		buff.items[strings.ToLower(k)] = v
	}

	return nil
}
