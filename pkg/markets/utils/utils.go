package utils

import (
	"io"

	"go.uber.org/zap"
)

func Dclose(c io.Closer, l *zap.Logger) {
	if err := c.Close(); err != nil {
		l.Error("Cant close response body", zap.Error(err))
	}
}
