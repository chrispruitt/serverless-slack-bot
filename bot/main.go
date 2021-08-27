package bot

import (
	"fmt"
	"os"
	"regexp"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var (
	SlackClient *slack.Client
	scripts     []Script
	BotName     string
)

type ScriptFunction func(*EventContext)

type Script struct {
	Name        string
	Matcher     Matcher
	Description string
	Function    ScriptFunction
}

func init() {

	BotName = getenv("BOT_NAME", "slackbot")

	SlackClient = slack.New(os.Getenv("SLACK_OAUTH_ACCESS_TOKEN"))
}

func RegisterScript(script Script) {
	scripts = append(scripts, script)
}

func HandleMentionEvent(event *slackevents.AppMentionEvent) {

	// Strip @bot-name out
	event.Text = stripBotName(event.Text)

	for _, script := range scripts {
		if match(script.Matcher.toRegex(), event.Text) {

			eventContext := &EventContext{
				SlackEvent: event,
			}

			eventContext.Arguments = script.Matcher.getArguments(event.Text)

			script.Function(eventContext)
			return
		}
	}

	PostMessage(event.Channel, "Sorry, I don't know that command.")
}

func PostMessage(channelID string, message string) (string, string, error) {
	if shellMode {
		fmt.Println(message)
		return "", "", nil
	} else {
		return SlackClient.PostMessage(channelID, slack.MsgOptionText(message, false))
	}
}

func match(matcher string, content string) bool {
	re := regexp.MustCompile(matcher)
	return re.MatchString(content)
}
