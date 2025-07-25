package aggregator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/utils"
	"go.uber.org/zap"
)

type Aggregator struct {
	markets map[string]markets.Market
	l       *zap.Logger
}

type MarketInfo struct {
	Market string
	Orders []markets.Pair
	URL    string
}

func New(markets map[string]markets.Market, l *zap.Logger) *Aggregator {
	return &Aggregator{
		markets: markets,
		l:       l,
	}
}

func (a *Aggregator) SearchAll(ctx context.Context, name string) []MarketInfo {
	var (
		wg = sync.WaitGroup{}
		mu = sync.Mutex{}
	)

	responses := make([]MarketInfo, 0)

	for marketName, svc := range a.markets {
		wg.Add(1)
		go func() {
			defer wg.Done()

			start := time.Now()
			defer u.RecordLatency(a.l, fmt.Sprintf("%s time to answer", marketName), start)

			res, err := svc.FindByHashName(ctx, name)
			url := svc.URL(name)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				responses = append(responses, MarketInfo{
					Market: marketName,
					Orders: nil,
					URL:    "",
				})
			} else {
				responses = append(responses, MarketInfo{
					Market: marketName,
					Orders: res,
					URL:    url,
				})
			}
		}()
	}

	wg.Wait()

	return responses
}
