package utils

import (
	"io"
	"time"

	"go.uber.org/zap"
)

func Dclose(c io.Closer, l *zap.Logger) {
	if err := c.Close(); err != nil {
		l.Error("Cant close response body", zap.Error(err))
	}
}

func RecordLatency(l *zap.Logger, msg string, start time.Time) {
	l.Debug(msg, zap.Duration("duration", time.Since(start)))
}
