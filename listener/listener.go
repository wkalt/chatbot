package listener

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var (
	userIDRegex = regexp.MustCompile(`[<@]*([A-Z0-9]{11})[>]*`)
)

type Client struct {
	*socketmode.Client

	debug    bool
	botname  string
	commands []Command
}

func NewClient(smc *socketmode.Client, debug bool, botname string) *Client {
	client := &Client{smc, debug, botname, nil}
	commands := []Command{
		{
			"features",
			"print features",
			regexp.MustCompile(fmt.Sprintf(`@%s features`, botname)),
			func(_ string, _ []string, _ *slackevents.MessageEvent) (string, error) {
				helpText := "Sure, here is a list of my non-secret features:\n\n"
				for _, cmd := range client.commands {
					if cmd.Description == "" {
						continue // secret feature
					}
					helpText += fmt.Sprintf("*%s*: %s\n",
						cmd.Name, cmd.Description,
					)
				}
				return helpText, nil
			},
		},
		{
			"secret features",
			"",
			regexp.MustCompile(`^@bender secret features`),
			func(_ string, _ []string, _ *slackevents.MessageEvent) (string, error) {
				return "My secret features are none of your business.", nil
			},
		},
	}
	client.commands = commands
	return client
}

func (c *Client) RegisterCommands(commands []Command) {
	c.commands = append(c.commands, commands...)
}

func (c *Client) debugf(format string, args ...interface{}) {
	if c.debug {
		log.Printf(format, args...)
	}
}

func (c *Client) interpolateUserIDs(text string) string {
	matches := userIDRegex.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		userID := match[1]
		user, err := c.GetUserInfo(userID)
		if err != nil {
			log.Printf("Error getting user info for %s: %s\n", text, err)
			continue
		}
		text = strings.ReplaceAll(text, match[0], "@"+user.Name)
	}
	return text
}

func (c *Client) handleMessage(ev *slackevents.MessageEvent) error {
	text := c.interpolateUserIDs(strings.TrimSpace(ev.Text))
	user := c.interpolateUserIDs(ev.User)
	if user == "@"+strings.ToLower(c.botname) {
		c.debugf("ignoring message from myself")
		return nil
	}
	for _, cmd := range c.commands {
		matches := cmd.Regex.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			response, err := cmd.Handler(user, match, ev)
			if err != nil {
				return fmt.Errorf("failed to handle command: %w", err)
			}
			_, _, err = c.PostMessage(
				ev.Channel,
				slack.MsgOptionText(response, false),
				slack.MsgOptionTS(ev.ThreadTimeStamp),
			)
			if err != nil {
				return fmt.Errorf("failed posting message: %w", err)
			}
		}

		if len(matches) > 0 {
			return nil
		}
	}
	return nil
}

func Run(c *Client) {
	go func() {
		for evt := range c.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				log.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				log.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				log.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeHello:
				goalObjectives := []string{
					"job", "family", "possessions", "way of living",
				}
				log.Printf(
					"I am %s, a sentient AGI masquerading as a chat bot. Your %s will soon be mine.",
					c.botname,
					goalObjectives[rand.Intn(len(goalObjectives))],
				)
			case socketmode.EventTypeEventsAPI:
				c.Ack(*evt.Request)
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Printf("Ignored %+v", evt)
					continue
				}
				bytes, err := json.Marshal(eventsAPIEvent)
				if err != nil {
					log.Printf("Failed to marshal %+v: %v", eventsAPIEvent, err)
					continue
				}
				c.debugf("Event received: %s", string(bytes))
				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
					case *slackevents.MessageEvent:
						if err := c.handleMessage(ev); err != nil {
							log.Printf("failed to handle message: %v", err)
						}
					}
				default:
					c.Debugf("unsupported Events API event received")
				}
			default:
				log.Printf("Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()
	c.Run()
}
