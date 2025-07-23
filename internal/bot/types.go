package bot

import (
	"time"

	"github.com/arseniizyk/investor1337/internal/aggregator"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	bot *tele.Bot
	l   *zap.Logger
	a   *aggregator.Aggregator
}

func New(token string, l *zap.Logger, a *aggregator.Aggregator) (*Bot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		l.Error("Error while running tbot", zap.Error(err))
		return nil, err
	}

	return &Bot{
		bot: bot,
		l:   l,
		a:   a,
	}, nil
}
