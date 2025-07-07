package aggregator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

type Aggregator struct {
	markets map[string]markets.Market
	l       *zap.Logger
}

func New(markets map[string]markets.Market, l *zap.Logger) *Aggregator {
	return &Aggregator{
		markets: markets,
		l:       l,
	}
}

func (a *Aggregator) SearchAll(ctx context.Context, name string) map[string]map[float64]int {
	var (
		wg = sync.WaitGroup{}
		mu = sync.Mutex{}
	)

	responses := make(map[string]map[float64]int, 0)

	for marketName, svc := range a.markets {
		wg.Add(1)
		go func() {
			defer wg.Done()

			start := time.Now()
			defer utils.RecordLatency(a.l, fmt.Sprintf("%s time to answer", marketName), start)

			res, err := svc.FindByHashName(ctx, name)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				responses[marketName] = make(map[float64]int, 0)
			} else {
				responses[marketName] = res
			}
		}()
	}

	wg.Wait()

	return responses
}
