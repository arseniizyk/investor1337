package bot

import (
	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type Tbot interface {
	Run() error
}

type tbot struct {
	token   string
	l       *zap.Logger
	markets map[string]markets.Market
}

func New(token string, l *zap.Logger, services map[string]markets.Market) Tbot {
	return tbot{
		token:   token,
		l:       l,
		markets: services,
	}
}
