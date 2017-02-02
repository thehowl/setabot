package bot

import (
	"fmt"
	"strconv"

	"gopkg.in/telegram-bot-api.v4"
)

var cities = map[string]string{
	"modena":             "mo",
	"piacenza":           "pc",
	"reggio emilia":      "re",
	"reggio nell'emilia": "re",
}

func (b *Bot) send(c tgbotapi.Chattable) {
	_, err := b.bot.Send(c)
	if err != nil {
		fmt.Println(err)
	}
}

// getCity retrieves the city of an user
func (b *Bot) getCity(u int) string {
	return b.Redis.Get("sbb:" + strconv.Itoa(u)).Val()
}

// setCity sets the city of an user
func (b *Bot) setCity(u int, city string) {
	b.Redis.Set("sbb:"+strconv.Itoa(u), city, 0)
}
