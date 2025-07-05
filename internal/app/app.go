package app

import (
	"github.com/arseniizyk/investor1337/internal/bot"
	"github.com/arseniizyk/investor1337/internal/config"
	"github.com/arseniizyk/investor1337/pkg/markets"
	"github.com/arseniizyk/investor1337/pkg/markets/aimmarket"
	"github.com/arseniizyk/investor1337/pkg/markets/csgomarket"
	"github.com/arseniizyk/investor1337/pkg/markets/csmoney"
	"github.com/arseniizyk/investor1337/pkg/markets/lisskins"
	"github.com/arseniizyk/investor1337/pkg/markets/steam"
	"go.uber.org/zap"
)

type App interface {
	Run() error
}

type app struct {
	c        config.Config
	l        *zap.Logger
	services map[string]markets.Market
}

func (a app) Run() error {
	tbot := bot.New(a.c.TelegramToken(), a.l, a.services)

	if err := tbot.Run(); err != nil {
		a.l.Error("Error while runnig tbot", zap.Error(err))
		return err
	}

	return nil
}

func New(cfg config.Config, l *zap.Logger) App {
	services := make(map[string]markets.Market, 1)

	cm := csgomarket.New(cfg.CsgoMarketToken(), l)
	csmoney := csmoney.New(l)
	ls := lisskins.New(cfg.LisSkinsToken(), l)
	services["CS GO Market"] = cm
	services["CS.Money"] = csmoney
	services["LIS-SKINS"] = ls

	steam, err := steam.New(l)
	if err != nil {
		l.Warn("Cant initialize steam", zap.Error(err))
	} else {
		services["STEAM"] = steam
	}

	am, err := aimmarket.New(l)
	if err != nil {
		l.Warn("Cant initialize aimmarket", zap.Error(err))
	} else {
		services["AIM.MARKET"] = am
	}

	return app{
		c:        cfg,
		l:        l,
		services: services,
	}
}
