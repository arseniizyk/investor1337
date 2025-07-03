package app

import (
	"github.com/arseniizyk/investor1337/internal/bot"
	"github.com/arseniizyk/investor1337/internal/config"
	"github.com/arseniizyk/investor1337/pkg/markets/csgomarket"
	"github.com/arseniizyk/investor1337/pkg/markets/csmoney"
	"github.com/arseniizyk/investor1337/pkg/markets/steam"
	"go.uber.org/zap"
)

type App interface {
	Run() error
}

type app struct {
	c config.Config
	l *zap.Logger
}

func (a app) Run() error {
	cm := csgomarket.New(a.c.CsgoMarketToken(), a.l)
	csmoney := csmoney.New(a.l)
	steam, err := steam.New(a.l)
	if err != nil {
		a.l.Error("Cant initialize steam")
	}

	tbot := bot.New(a.c.TelegramToken(), a.l, cm, csmoney, steam)

	if err := tbot.Run(); err != nil {
		a.l.Error("Error while runnig tbot", zap.Error(err))
		return err
	}

	return nil
}

func New(cfg config.Config, l *zap.Logger) App {
	return app{cfg, l}
}
