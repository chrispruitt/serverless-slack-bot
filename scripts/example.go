package scripts

import (
	"fmt"
	"strings"

	"github.com/chrispruitt/serverless-slack-bot/bot"
)

func init() {
	// Simple script
	bot.RegisterScript(bot.Script{
		Name:        "lulz",
		Matcher:     "lulz",
		Description: "lulz",
		Function: func(context *bot.EventContext) {
			bot.PostMessage(context.SlackEvent.Channel, "lol")
		},
	})

	// Script with parameters
	bot.RegisterScript(bot.Script{
		Name:        "echo",
		Matcher:     "echo <message>",
		Description: "Echo a message",
		Function: func(context *bot.EventContext) {
			message := context.Arguments["message"]
			bot.PostMessage(context.SlackEvent.Channel, fmt.Sprintf("You said, \"%s\"", message))
		},
	})

	// Script with some custom parameter syntax
	bot.RegisterScript(bot.Script{
		Name:        "ship",
		Matcher:     `ship <app> to <env>`,
		Description: "Usage: 'ship app1@v1.0.0 to dev' or 'ship app1@v1.0.0 app2@v1.0.0 to dev",
		Function: func(context *bot.EventContext) {
			// TODO Validation
			env := context.Arguments["env"]
			apps := strings.Split(context.Arguments["app"], " ")

			for _, a := range apps {
				app := strings.Split(a, "@")
				bot.PostMessage(context.SlackEvent.Channel, fmt.Sprintf("Shipping App: %s Version: %s to %s", app[0], app[1], env))
			}
		},
	})
}
