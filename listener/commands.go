package listener

import (
	"regexp"

	"github.com/slack-go/slack/slackevents"
)

var motdCache = map[string]string{}

type CommandHandler func(user string, matches []string, ev *slackevents.MessageEvent) (string, error)

type Command struct {
	Name        string
	Description string
	Regex       *regexp.Regexp
	Handler     CommandHandler
}

func MakeCommand(
	name, description string, regex *regexp.Regexp, handler CommandHandler) Command {
	return Command{name, description, regex, handler}
}
