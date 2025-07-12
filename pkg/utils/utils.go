package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type Response interface {
	LenData() int
	Cursor() string
}

var (
	ErrRequest = errors.New("cant make request")
)

func Dclose(c io.Closer, l *zap.Logger) {
	if err := c.Close(); err != nil {
		l.Error("Cant close response body", zap.Error(err))
	}
}

func RecordLatency(l *zap.Logger, msg string, start time.Time) {
	l.Debug(msg, zap.Duration("duration", time.Since(start)))
}

func DoJSONRequest[T any](ctx context.Context, client *http.Client, req *http.Request, logger *zap.Logger) (T, error) {
	var zero T

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("HTTP request failed", zap.Error(err))
		return zero, err
	}

	defer Dclose(resp.Body, logger)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		logger.Warn("Bad status_code", zap.Int("status_code", resp.StatusCode))
		return zero, fmt.Errorf("status %d", resp.StatusCode)
	}

	var res T

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		logger.Error("Cant decode response", zap.Error(err))
		return zero, err
	}

	return res, nil
}

func FetchWithCursor[T Response](
	ctx context.Context, client *http.Client, l *zap.Logger,
	name, marketName string,
	maxPages int,
	countInMap func(m map[float64]int, r *T),
	buildRequest func(cursor, name string) (*http.Request, error)) (map[float64]int, error) {

	countMap := make(map[float64]int, markets.MaxOutputs)
	cursor := ""

	for i := range maxPages {
		if len(countMap) == markets.MaxOutputs {
			break
		}

		req, err := buildRequest(cursor, name)
		if err != nil {
			return nil, err
		}

		r, err := DoJSONRequest[T](ctx, client, req, l)
		if err != nil || r.LenData() == 0 {
			l.Warn("Response error",
				zap.String("market", marketName),
				zap.String("name", name),
				zap.Error(err),
			)
			return nil, err
		}

		countInMap(countMap, &r)

		l.Debug("Fetched page",
			zap.String("market", marketName),
			zap.Int("page", i),
			zap.String("name", name),
			zap.String("cursor", cursor),
		)

		if r.Cursor() == "" {
			break
		}

		cursor = r.Cursor()

		if r.LenData() == 0 {
			break
		}
	}

	return countMap, nil
}

func SortPairs(pairs []markets.Pair) {
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].Price < pairs[j].Price
	})
}

func SinglePair(price float64, count int) []markets.Pair {
	return []markets.Pair{{
		Price:    price,
		Quantity: count,
	}}
}

func PairsFromMap(m map[float64]int) []markets.Pair {
	result := make([]markets.Pair, 0, len(m))

	for price, quantity := range m {
		result = append(result, markets.Pair{
			Price:    price,
			Quantity: quantity,
		})
	}

	SortPairs(result)

	return result
}
