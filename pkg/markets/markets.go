package markets

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/arseniizyk/investor1337/pkg/utils"
	"go.uber.org/zap"
)

const MaxOutputs = 4

var (
	ErrNoOffers      = errors.New("bad request, no offers")
	ErrFormatFailed  = errors.New("cant format JSON response")
	ErrBadResponse   = errors.New("bad response")
	ErrRequestFailed = errors.New("failed to build HTTP request")
	ErrBadStatusCode = errors.New("unexpected HTTP status code")
	ErrDecodeJSON    = errors.New("failed to decode JSON response")
	ErrEmptyResponse = errors.New("empty or invalid response")
)

type Market interface {
	FindByHashName(ctx context.Context, name string) ([]Pair, error)
	URL(name string) string
}

type Response interface {
	LenData() int
	Cursor() string
}

func DoJSONRequest[T any](ctx context.Context, client *http.Client, req *http.Request, logger *zap.Logger) (T, error) {
	var zero T

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("HTTP request failed", zap.Error(err))
		return zero, ErrRequestFailed
	}

	defer utils.Dclose(resp.Body, logger)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		logger.Warn("Bad status_code", zap.Int("status_code", resp.StatusCode))
		return zero, ErrBadStatusCode
	}

	var res T

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		logger.Error("Cant decode response", zap.Error(err))
		return zero, ErrDecodeJSON
	}

	return res, nil
}

func FetchWithCursor[T Response](
	ctx context.Context, client *http.Client, l *zap.Logger,
	name, marketName string,
	maxPages int,
	countInMap func(m map[float64]int, r *T),
	buildRequest func(cursor, name string) (*http.Request, error)) (map[float64]int, error) {

	countMap := make(map[float64]int, MaxOutputs)
	cursor := ""

	for i := range maxPages {
		if len(countMap) == MaxOutputs {
			break
		}

		req, err := buildRequest(cursor, name)
		if err != nil {
			l.Error("Cant build request",
				zap.String("market", marketName),
				zap.String("name", name),
				zap.Error(err),
			)
			return nil, ErrRequestFailed
		}

		r, err := DoJSONRequest[T](ctx, client, req, l)
		if err != nil || r.LenData() == 0 {
			l.Warn("Response error",
				zap.String("market", marketName),
				zap.String("name", name),
				zap.Error(err),
			)
			return nil, ErrEmptyResponse
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
