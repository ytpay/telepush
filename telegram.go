package main

import (
	"io"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var bot *Telegram

type Telegram struct {
	bot *tb.Bot
}

func NewTelegram(api, token string) (*Telegram, error) {
	bot, err := tb.NewBot(tb.Settings{
		URL:    api,
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	})
	if err != nil {
		return nil, err
	} else {
		return &Telegram{bot: bot}, nil
	}
}

func (tg *Telegram) ID() {
	tg.bot.Handle("/id", func(m *tb.Message) {
		if !m.Private() {
			_, _ = tg.bot.Send(m.Sender, m.Sender.ID)
		} else {
			_, _ = tg.bot.Send(m.Chat, m.Chat.ID)
		}
	})
	tg.bot.Start()
}

func (tg *Telegram) SendMessage(msg string, to int64, markdown bool) error {
	opt := &tb.SendOptions{}
	if markdown {
		opt.ParseMode = tb.ModeMarkdown
	}

	_, err := tg.bot.Send(tb.ChatID(to), msg, opt)
	return err
}

func (tg *Telegram) SendFile(file io.Reader, fileName, mime, caption string, to int64) error {
	_, err := tg.bot.Send(tb.ChatID(to), &tb.Document{
		File:     tb.File{FileReader: file},
		Caption:  caption,
		MIME:     mime,
		FileName: fileName,
	})
	return err
}

func (tg *Telegram) SendImage(image io.Reader, caption string, to int64) error {
	_, err := tg.bot.Send(tb.ChatID(to), &tb.Photo{
		File:    tb.File{FileReader: image},
		Caption: caption,
	})
	return err
}
