package bot

import (
	"fmt"
	"sort"
	"strings"
	"time"

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
	skin := c.Text()

	var (
		cmRes      map[float64]int
		csmoneyRes map[float64]int
		steamRes   map[float64]int
		err        error
	)

	if cmRes, err = t.cstm.FindByHashName(skin); err != nil {
		t.l.Warn("CS Market error", zap.Error(err))
		cmRes = map[float64]int{}
	}

	if csmoneyRes, err = t.csmoney.FindByHashName(skin); err != nil {
		t.l.Warn("CS Money error", zap.Error(err))
		csmoneyRes = map[float64]int{}
	}

	if steamRes, err = t.steam.FindByHashName(strings.ToLower(skin)); err != nil {
		t.l.Warn("Steam error", zap.Error(err))
		steamRes = map[float64]int{}
	}

	msg := format(cmRes, "CS Market") +
		format(csmoneyRes, "CSmoney") +
		format(steamRes, "Steam")

	return c.Send(msg)
}

func format(res map[float64]int, market string) string {
	if len(res) == 0 {
		return fmt.Sprintf("%s не найдено\n\n", market)
	}

	keys := make([]float64, 0, len(res))
	for k := range res {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	result := fmt.Sprintf("%s\n", market)
	for _, k := range keys {
		result += fmt.Sprintf("Price: $%.2f | %d\n", k, res[k])
	}
	result += "\n"
	return result
}
