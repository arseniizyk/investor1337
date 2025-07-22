package bot

import (
	"context"
	"time"

	u "github.com/arseniizyk/investor1337/pkg/utils"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (b *Bot) Run() error {
	b.l.Info("Bot running", zap.String("username", b.bot.Me.Username))

	b.bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Type skin name (e.g Fracture Case)")
	})
	b.bot.Handle(tele.OnText, b.findByName)

	b.bot.Start()

	return nil
}

func (b *Bot) findByName(c tele.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	defer u.RecordLatency(b.l, "findByName time to answer", start)
	responses := b.a.SearchAll(ctx, c.Text())

	msg := format(responses)
	return c.Send(msg, tele.ModeMarkdown, tele.NoPreview)
}
