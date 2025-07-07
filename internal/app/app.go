package app

import (
	"net/http"
	"time"

	"github.com/arseniizyk/investor1337/internal/aggregator"
	"github.com/arseniizyk/investor1337/internal/bot"
	"github.com/arseniizyk/investor1337/internal/config"
	"github.com/arseniizyk/investor1337/pkg/markets"
	"github.com/arseniizyk/investor1337/pkg/markets/aimmarket"
	"github.com/arseniizyk/investor1337/pkg/markets/csfloat"
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
	c config.Config
	l *zap.Logger
	a *aggregator.Aggregator
}

func (a app) Run() error {
	tbot := bot.New(a.c.TelegramToken(), a.l, a.a)

	if err := tbot.Run(); err != nil {
		a.l.Error("Error while runnig tbot", zap.Error(err))
		return err
	}

	return nil
}

func New(cfg config.Config, l *zap.Logger) App {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	markets := map[string]markets.Market{
		"CS GO Market": csgomarket.New(client, cfg.CsgoMarketToken(), l),
		"CS.Money":     csmoney.New(client, l),
		"LIS-SKINS":    lisskins.New(client, cfg.LisSkinsToken(), l),
		"CSFloat":      csfloat.New(client, cfg.CsfloatCookie(), l),
	}

	steam, err := steam.New(client, l)
	if err != nil {
		l.Warn("Cant initialize steam", zap.Error(err))
	} else {
		markets["STEAM"] = steam
	}

	am, err := aimmarket.New(client, l)
	if err != nil {
		l.Warn("Cant initialize aimmarket", zap.Error(err))
	} else {
		markets["AIM.MARKET"] = am
	}

	agg := aggregator.New(markets, l)

	return app{
		c: cfg,
		l: l,
		a: agg,
	}
}
