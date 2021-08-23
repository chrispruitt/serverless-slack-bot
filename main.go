package main

import (
	"os"

	"github.com/chrispruitt/serverless-slack-bot/bot"
	_ "github.com/chrispruitt/serverless-slack-bot/scripts"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "shell" {
		bot.Shell()
	} else {
		bot.Start()
	}
}
