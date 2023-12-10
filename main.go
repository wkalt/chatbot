package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack/socketmode"
	"github.com/wkalt/chatbot/external"
	"github.com/wkalt/chatbot/listener"

	"github.com/slack-go/slack"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No envfile found")
	}
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must be set.\n")
		os.Exit(1)
	}
	if !strings.HasPrefix(appToken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	botName := os.Getenv("BOT_NAME")
	if botName == "" {
		fmt.Fprintf(os.Stderr, "BOT_NAME must be set.\n")
		os.Exit(1)
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must be set.\n")
		os.Exit(1)
	}
	if !strings.HasPrefix(botToken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}
	debug := os.Getenv("DEBUG") == "true"

	api := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
	)
	external.Init()
	client := listener.NewClient(socketmode.New(api), debug, botName)
	client.RegisterCommands(external.Commands)
	listener.Run(client)
}
