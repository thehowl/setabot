package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/thehowl/setabusbot/bot"
	"github.com/thehowl/setabusbot/providers/scraper"
	redis "gopkg.in/redis.v5"
)

func main() {
	i, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	b := &bot.Bot{
		Redis: redis.NewClient(&redis.Options{
			Network:  getenv("REDIS_NETWORK", "tcp"),
			Addr:     getenv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       i,
		}),
		BotToken: os.Getenv("BOT_TOKEN"),
		AS:       &scraper.Scraper{},
	}
	fmt.Println(b.Start())
	os.Exit(1)
}

func getenv(v string, def string) string {
	v = os.Getenv(v)
	if v == "" {
		return def
	}
	return v
}
