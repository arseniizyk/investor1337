package bot

import (
	"context"
	"time"

	u "github.com/arseniizyk/investor1337/pkg/utils"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (t tbot) Run() error {
	pref := tele.Settings{
		Token:  t.token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		t.l.Error("Error while running tbot", zap.Error(err))
		return err
	}

	t.l.Info("Bot running", zap.String("username", b.Me.Username))

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Type skin name (e.g Fracture Case)")
	})

	b.Handle(tele.OnText, t.findByName)

	b.Start()

	return nil
}

func (t tbot) findByName(c tele.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	start := time.Now()
	defer u.RecordLatency(t.l, "findByName time to answer", start)
	responses := t.a.SearchAll(ctx, c.Text())

	msg := format(responses)
	return c.Send(msg, tele.ModeMarkdown, tele.NoPreview)
}
