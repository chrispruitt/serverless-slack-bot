**Description**

This is a simple hubot like slack bot with serverless in mind. In short it is an API to support the mention event of slack apps and parses the event message to run your custom scripts. It was written with serverless in mind so you can host it using lambda and api-gatway with minimal set up.



**Setup**

1. Create your slack app and get your bot tokens
2. Create a main.go file and scripts/ folder
3. Setup lambda function with whatever permissions required by your custom scripts
4. Setup API Gateway with a endpoint `/{proxy+}` pointing to your lambda function
5. Deploy your code to your lambda function
6. Configure your slack app to point to your api gateway



main.go example

```go
package main

import (
	"github.com/chrispruitt/ssbot/bot"

	_ "<yourModuleName>/scripts"
)

func main() {
	bot.Start()
}
```

scripts/exampleScript.go

```go
package scripts

import (
	"fmt"
	"regexp"

	"github.com/chrispruitt/ssbot/bot"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func init() {
	bot.RegisterScript(bot.Script{
		Name:               "Echo",
		Matcher:            "(?i)^echo.*",
		Description:        "Echo a message",
		CommandDescription: "echo <message>",
		Function: func(event *slackevents.AppMentionEvent) {
			re := regexp.MustCompile(`echo *`)
			text := re.ReplaceAllString(event.Text, "")

			bot.SlackClient.PostMessage(event.Channel, slack.MsgOptionText(fmt.Sprintf("You said, \"%s\"", text), false))
		},
	})
}

```



**Serve API Locally**

```bash
export SLACK_OAUTH_ACCESS_TOKEN=yourslackoauthaccecctoken
export SLACK_VERIFICATION_TOKEN=yourslackapiverificationtoken

go run main.go
```



**In Slack**

Add your slack bot to a channel.

Then, execute a script via slack by mentioning your slack bot followed by a scripts command.



**Roadmap**

- Provide terraform module for quick setup of lambda and api-gateway
- Update readme with a "how to" to set up slack bot or publish one
- Add script authorization via roles
- Give ssbot a brain via dynamodb or s3 json file