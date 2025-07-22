package app

import (
	"net/http"
	"time"

	"github.com/arseniizyk/investor1337/internal/aggregator"
	"github.com/arseniizyk/investor1337/internal/bot"
	"github.com/arseniizyk/investor1337/internal/config"
	"github.com/arseniizyk/investor1337/pkg/markets"
	"github.com/arseniizyk/investor1337/pkg/markets/aimmarket"
	"github.com/arseniizyk/investor1337/pkg/markets/buff163"
	"github.com/arseniizyk/investor1337/pkg/markets/csfloat"
	"github.com/arseniizyk/investor1337/pkg/markets/csgomarket"
	"github.com/arseniizyk/investor1337/pkg/markets/csmoney"
	"github.com/arseniizyk/investor1337/pkg/markets/lisskins"
	"github.com/arseniizyk/investor1337/pkg/markets/steam"
	"go.uber.org/zap"
)

type App struct {
	c config.Config
	l *zap.Logger
	a *aggregator.Aggregator
}

func (a *App) Run() error {
	tbot, err := bot.New(a.c.TelegramToken(), a.l, a.a)
	if err != nil {
		return err
	}

	if err := tbot.Run(); err != nil {
		a.l.Error("Error while running tbot", zap.Error(err))
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

	buff, err := buff163.New(client, l)
	if err != nil {
		l.Warn("Cant initialize buff163", zap.Error(err))
	} else {
		markets["BUFF163"] = buff
	}

	am, err := aimmarket.New(client, l)
	if err != nil {
		l.Warn("Cant initialize aimmarket", zap.Error(err))
	} else {
		markets["AIM.MARKET"] = am
	}

	agg := aggregator.New(markets, l)

	return App{
		c: cfg,
		l: l,
		a: agg,
	}
}
