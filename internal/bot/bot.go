package bot

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"github.com/arseniizyk/investor1337/pkg/markets/utils"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (t tbot) Run() error {
	pref := tele.Settings{
		Token:  t.token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		t.l.Error("Error while running tbot", zap.Error(err))
		return err
	}

	t.l.Info("Bot running", zap.String("username", b.Me.Username))

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Type skin name (e.g Fracture Case)")
	})

	b.Handle(tele.OnText, t.findByName)

	b.Start()

	return nil
}

func (t tbot) findByName(c tele.Context) error {
	var (
		wg = sync.WaitGroup{}
		mu = sync.Mutex{}
	)

	skin := c.Text()
	responses := make(map[string]map[float64]int, 0)

	start := time.Now()
	defer utils.RecordLatency(t.l, "findByName time to answer", start)

	for name, svc := range t.markets {
		wg.Add(1)
		go func(name string, svc markets.Market) {
			start := time.Now()
			defer utils.RecordLatency(t.l, fmt.Sprintf("%s time to answer", name), start)

			defer wg.Done()
			res, err := svc.FindByHashName(skin) // TODO provide context
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				responses[name] = make(map[float64]int, 0)
			} else {
				responses[name] = res
			}
		}(name, svc)
	}

	wg.Wait()

	msg := format(responses)

	return c.Send(msg)
}

func format(res map[string]map[float64]int) string {
	var result strings.Builder

	for market, offers := range res {
		if len(offers) == 0 {
			result.WriteString(fmt.Sprintf("%s не найдено предложений\n\n", market))
			continue
		}

		keys := make([]float64, 0, len(offers))
		for k := range offers {
			keys = append(keys, k)
		}
		sort.Float64s(keys)

		result.WriteString(fmt.Sprintf("%s\n", market))
		for _, k := range keys {
			result.WriteString(fmt.Sprintf("Price: $%.2f | %d\n", k, offers[k]))
		}
		result.WriteString("\n")
	}

	return result.String()
}
