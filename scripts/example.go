package scripts

import (
	"fmt"
	"regexp"

	"github.com/chrispruitt/serverless-slack-bot/bot"

	"github.com/slack-go/slack/slackevents"
)

func init() {
	bot.RegisterScript(bot.Script{
		Name:               "lulz",
		Matcher:            "(?i)^lulz$",
		Description:        "lulz",
		CommandDescription: "lulz",
		Function: func(event *slackevents.AppMentionEvent) {

			bot.PostMessage(event.Channel, "lol")
		},
	})

	bot.RegisterScript(bot.Script{
		Name:               "Echo",
		Matcher:            "(?i)^echo.*",
		Description:        "Echo a message",
		CommandDescription: "echo <message>",
		Function: func(event *slackevents.AppMentionEvent) {
			re := regexp.MustCompile(`echo *`)
			text := re.ReplaceAllString(event.Text, "")

			bot.PostMessage(event.Channel, fmt.Sprintf("You said, \"%s\"", text))
		},
	})
}
