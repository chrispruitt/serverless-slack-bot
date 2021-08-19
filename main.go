package main

import (
	"github.com/chrispruitt/serverless-slack-bot/bot"

	_ "github.com/chrispruitt/serverless-slack-bot/scripts"
)

func main() {
	bot.Start()
}
