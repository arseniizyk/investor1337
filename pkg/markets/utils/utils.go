package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
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

func SortPairs(pairs []markets.Pair) {
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].Price < pairs[j].Price
	})
}
