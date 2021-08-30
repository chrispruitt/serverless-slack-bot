package bot

import "github.com/slack-go/slack/slackevents"

type EventContext struct {
	Arguments  map[string]string
	SlackEvent *slackevents.AppMentionEvent
}
