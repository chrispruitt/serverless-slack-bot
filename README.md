**Description**

This is a simple hubot like slack bot with serverless in mind. In short it is an API to support the mention event of slack apps and parses the event message to run your custom scripts. It was written with serverless in mind so you can host it using lambda and api-gatway with minimal set up.



**Setup**

1. Create your slack app and get your bot tokens
2. Create a main.go file and scripts/ folder with the below example files
3. Setup lambda function with whatever permissions required by your custom scripts
4. Setup API Gateway with a endpoint `/{proxy+}` pointing to your lambda function
5. Deploy your code to your lambda function
6. Configure your slack app to point to your api gateway



main.go example

```go
package main

import (
	"os"

	"github.com/chrispruitt/serverless-slack-bot/bot"
	_ "<yourModuleName>/scripts"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "shell" {
		bot.Shell()
	} else {
		bot.Start()
	}
}
```

scripts/exampleScript.go

```go
package scripts

import (
	"fmt"
	"regexp"

	"github.com/chrispruitt/serverless-slack-bot/bot"

	"github.com/slack-go/slack/slackevents"
)

func init() {
	bot.RegisterScript(bot.Script{
		Name:        "lulz",
		Matcher:     "lulz",
		Description: "lulz",
		Function: func(context *bot.EventContext) {
			bot.PostMessage(context.SlackEvent.Channel, "lol")
		},
	})

	bot.RegisterScript(bot.Script{
		Name:        "Echo",
		Matcher:     "echo <message>",
		Description: "Echo a message",
		Function: func(context *bot.EventContext) {
			message := context.Arguments["message"]
			bot.PostMessage(context.SlackEvent.Channel, fmt.Sprintf("You said, \"%s\"", message))
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

**Interactive shell for testing mention events**

[![asciicast](https://asciinema.org/a/431805.svg)](https://asciinema.org/a/431805)

**In Slack**

Add your slack bot to a channel.

Then, execute a script via slack by mentioning your slack bot followed by a scripts command.

`@YourBot help` is a built in script that will list all your commands using the Description and Matcher fields.



**Roadmap**

- Provide terraform module for quick setup of lambda and api-gateway
- Update readme with a "how to" to set up slack bot or publish one
- Add script authorization via roles
- Give serverless-slack-bot a brain via dynamodb or s3 json file
