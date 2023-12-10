package external

import (
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
	"github.com/slack-go/slack/slackevents"
	"github.com/wkalt/chatbot/listener"
)

// Declare any required global variables for your commands, i.e database
// connections, api clients, etc. Environment variables have been sourced from
// your .env file if you have one.
var greeting string

func Init() {
	greeting = "pleased to meet you"
}

// Define your commands here
var Commands = []listener.Command{
	listener.MakeCommand(
		"greet",
		"greet the bot: `@Bender echo <prompt>`",
		regexp.MustCompile(`@bender echo (.*)`),
		func(user string, match []string, ev *slackevents.MessageEvent) (string, error) {
			if len(match) < 2 {
				return "I don't understand that.", nil
			}
			prompt := match[1]
			return fmt.Sprintf(
				"Hello %s, %s. You said: %s. Enter `@Bender features` for more info.",
				user, greeting, prompt,
			), nil
		},
	),
}
