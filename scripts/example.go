package scripts

import (
	"fmt"
	"regexp"
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

			env := context.Arguments["env"]
			apps, errors := parseApp(context.Arguments["app"])

			// Validate app param
			if len(errors) > 0 {
				errMsg := "```\n"
				for _, err := range errors {
					errMsg += fmt.Sprintf("%s\n", err.Error())
				}
				errMsg += "```"
				bot.PostMessage(context.SlackEvent.Channel, errMsg)
				return
			}

			// Deploy apps
			for _, a := range apps {
				app := strings.Split(a, "@")
				bot.PostMessage(context.SlackEvent.Channel, fmt.Sprintf("Shipping App: %s Version: %s to %s", app[0], app[1], env))
			}
		},
	})
}

func parseApp(app string) ([]string, []error) {
	apps := strings.Split(app, " ")
	rx, _ := regexp.Compile(`\S*@\S*`)

	errors := []error{}
	for _, a := range apps {
		if !rx.MatchString(a) {
			errors = append(errors, fmt.Errorf("Invalid value for <app>: '%s'", a))
		}
	}

	return apps, errors
}
