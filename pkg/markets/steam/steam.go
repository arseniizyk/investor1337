package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

func (s steam) FindByHashName(name string) (map[float64]int, error) {
	url := fmt.Sprintf("https://steamcommunity.com/market/itemordershistogram?norender=1&language=english&currency=1&item_nameid=%d", s.data[name])

	resp, err := http.Get(url)
	if err != nil {
		s.l.Error("cant request steam", zap.String("name", name), zap.Error(err))
		return nil, err
	}
	defer utils.Dclose(resp.Body, s.l)

	switch resp.StatusCode {
	case http.StatusOK:
		var res Response
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			s.l.Error("Cant decode json", zap.String("name", name), zap.Error(err))
			return nil, err
		}

		results, err := format(res)
		if err != nil {
			s.l.Error("Cant format response from steam", zap.String("name", name), zap.Error(err))
			return nil, err
		}

		return results, nil

	default:
		s.l.Warn("Unknown status code from steam", zap.Int("status_code", resp.StatusCode))
		return nil, err
	}
}

func format(r Response) (map[float64]int, error) {
	results := make(map[float64]int)

	for i, orders := range r.SellOrderTable {
		if i == 3 {
			break
		}

		re := regexp.MustCompile(`\d+\.?\d*`)
		comp := re.FindString(orders.Price)

		price, err := strconv.ParseFloat(comp, 64)
		if err != nil {
			return nil, err
		}

		quantity, err := strconv.Atoi(strings.ReplaceAll(orders.Quantity, ",", ""))
		if err != nil {
			return nil, err
		}

		results[price] = quantity
	}

	return results, nil
}

func (s *steam) loadNameIds() error {
	file, err := os.ReadFile("../cs2ids.json")
	if err != nil {
		s.l.Error("Cant load cs2 ids from json", zap.Error(err))
		return err
	}

	data := make(map[string]int)

	if err := json.Unmarshal(file, &data); err != nil {
		s.l.Error("Cant unmarshal cs2 ids", zap.Error(err))
		return err
	}

	s.data = make(map[string]int)
	for k, v := range data {
		lower := strings.ToLower(k)
		s.data[lower] = v
	}

	return nil
}
