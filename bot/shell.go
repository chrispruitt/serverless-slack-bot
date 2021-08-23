package bot

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattes/go-asciibot"
	"github.com/slack-go/slack/slackevents"

	"github.com/common-nighthawk/go-figure"
)

var (
	shellMode bool
)

func init() {
	shellMode = false
}

func Shell() {
	Banner()
	shellMode = true
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("slackbot> ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = runCommand(cmdString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func Banner() {
	myFigure := figure.NewColorFigure("Slack bot", "", "green", true)
	myFigure.Print()
	fmt.Println(asciibot.Random())
	fmt.Println("Welcome to the slackbot shell! Type 'help' for help, 'exit' to exit.")
}

func runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	switch arrCommandStr[0] {
	case "exit":
		os.Exit(0)
	default:
		event := &slackevents.AppMentionEvent{
			Text:    commandStr,
			Channel: "shell",
		}
		HandleMentionEvent(event)
		return nil
	}
	cmd := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
