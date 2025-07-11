package aggregator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/arseniizyk/investor1337/pkg/markets"
	u "github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
)

type Output struct {
	Market string
	Orders []markets.Pair
}

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

func (a *Aggregator) SearchAll(ctx context.Context, name string) []Output {
	var (
		wg = sync.WaitGroup{}
		mu = sync.Mutex{}
	)

	responses := make([]Output, 0)

	for marketName, svc := range a.markets {
		wg.Add(1)
		go func() {
			defer wg.Done()

			start := time.Now()
			defer u.RecordLatency(a.l, fmt.Sprintf("%s time to answer", marketName), start)

			res, err := svc.FindByHashName(ctx, name)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				a.l.Sugar().Error("error in FindByHashName:", marketName)
			} else {
				responses = append(responses, Output{Market: marketName, Orders: res})
			}
		}()
	}

	wg.Wait()

	return responses
}
