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
	cstm    markets.Market
	csmoney markets.Market
	steam   markets.Market
}

func New(token string, l *zap.Logger, cstm, csmoney, steam markets.Market) Tbot {
	return tbot{token, l, cstm, csmoney, steam}
}
