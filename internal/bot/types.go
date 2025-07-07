package bot

import (
	"github.com/arseniizyk/investor1337/internal/aggregator"
	"go.uber.org/zap"
)

type Tbot interface {
	Run() error
}

type tbot struct {
	token string
	l     *zap.Logger
	a     *aggregator.Aggregator
}

func New(token string, l *zap.Logger, a *aggregator.Aggregator) Tbot {
	return tbot{
		token: token,
		l:     l,
		a:     a,
	}
}
