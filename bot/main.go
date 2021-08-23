package bot

import (
	"fmt"
	"os"
	"regexp"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var SlackClient *slack.Client
var scripts []Script

type ScriptFunction func(*slackevents.AppMentionEvent)

type Script struct {
	Name               string
	Matcher            string
	Description        string
	CommandDescription string
	Function           ScriptFunction
}

func init() {

	SlackClient = slack.New(os.Getenv("SLACK_OAUTH_ACCESS_TOKEN"))

	RegisterScript(Script{
		Name:               "Help",
		Matcher:            "(?i)^help$",
		Description:        "show description for all commands",
		CommandDescription: "help",
		Function:           helpScriptFunc,
	})
}

func RegisterScript(script Script) {
	scripts = append(scripts, script)
}

func HandleMentionEvent(event *slackevents.AppMentionEvent) {

	// Strip @bot-name out
	re := regexp.MustCompile(`^<@.*> *`)
	event.Text = re.ReplaceAllString(event.Text, "")

	for _, script := range scripts {
		if match(script.Matcher, event.Text) {
			script.Function(event)
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

func helpScriptFunc(event *slackevents.AppMentionEvent) {
	helpMsg := "Prefix @<yourBotUser> to any command you would like to execute. \n\n"
	for i, script := range scripts {
		if i != 0 {
			helpMsg += "\n"
		}
		if script.CommandDescription != "" {
			helpMsg += "@slackbot " + script.CommandDescription
			if script.Description != "" {
				helpMsg += fmt.Sprintf(" - %s", script.Description)
			}
		} else {
			helpMsg += fmt.Sprintf("Missing help command description for %s", script.Name)
		}
	}
	PostMessage(event.Channel, fmt.Sprintf("```%s```", helpMsg))
}
