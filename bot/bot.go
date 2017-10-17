package bot

import (
	"fmt"
	"log"

	"strings"

	"github.com/thehowl/setabot/services"
	"gopkg.in/redis.v5"
	"gopkg.in/telegram-bot-api.v4"
)

// Bot is an instance of setabot.
type Bot struct {
	Redis    *redis.Client
	BotToken string
	bot      *tgbotapi.BotAPI
	AS       services.ArrivalsService
	commands map[string]func(u tgbotapi.Update)
}

// Start begins taking updates from Telegram's API
func (b *Bot) Start() error {
	b.commands = map[string]func(u tgbotapi.Update){
		"/start":  b.start,
		"/qm":     b.qm,
		"/search": b.search,
		"sono di": b.imFrom,
	}

	var err error
	b.bot, err = tgbotapi.NewBotAPI(b.BotToken)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)

	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message != nil {
			fmt.Println("@"+update.Message.From.UserName, "-", update.Message.Text)
			go b.handleUpdate(update)
		}
	}

	return nil
}

func (b *Bot) handleUpdate(u tgbotapi.Update) {
	// we are adding a space to the text because we nevertheless want to match
	// commands like /start (with no arguments).
	txt := strings.ToLower(u.Message.Text + " ")
	for cname, han := range b.commands {
		// we want to match /qm Galilei. We don't want to match /qmGalilei.
		cname += " "
		if strings.HasPrefix(txt, cname) {
			u.Message.Text = strings.TrimSpace(strings.TrimPrefix(txt, cname))
			han(u)
			return
		}
	}
}
