package bot

import (
	"fmt"
	"strings"
)

func init() {
	RegisterScript(Script{
		Name:        "Help",
		Matcher:     "help ?<filter>",
		Description: "show description for all commands",
		Function:    helpScriptFunc,
	})
}

func helpScriptFunc(context *EventContext) {
	filter := context.Arguments["filter"]
	helpMsg := ""
	for _, script := range scripts {
		if script.Matcher != "" && (strings.Contains(string(script.Matcher), filter) || filter == "") {
			if helpMsg != "" {
				helpMsg += "\n"
			}
			helpMsg = fmt.Sprintf("%s%s", helpMsg, string(script.Matcher))
			if script.Description != "" {
				helpMsg += fmt.Sprintf(" - %s", script.Description)
			} else {
				helpMsg += fmt.Sprintf("Missing help command description for %s", script.Name)
			}
		}
	}
	message := fmt.Sprintf("```\n%s\n```", helpMsg)
	if shellMode {
		message = helpMsg
	}
	PostMessage(context.SlackEvent.Channel, message)
}
