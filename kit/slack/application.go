package slack

import (
	"fmt"
	"github.com/gonzispina/channeled/kit/context"
	"github.com/gonzispina/channeled/kit/logs"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"strings"
	"time"
)

// NewApplication constructor
func NewApplication(c *Config, l logs.Logger) *ApplicationController {
	if c == nil {
		panic("Config must be initialized")
	}

	if l == nil {
		panic("Config must be initialized")
	}

	if c.AppName == "" {
		panic("SLACK_APP_TOKEN must be set")
	}

	if !strings.HasPrefix(c.AppToken, "xapp-") {
		panic("SLACK_APP_TOKEN must have the prefix \"xapp-\"")
	}

	if c.WorkspaceToken == "" {
		panic("SLACK_BOT_TOKEN must be set")
	}

	if !strings.HasPrefix(c.WorkspaceToken, "xoxb-") {
		panic("SLACK_BOT_TOKEN must have the prefix \"xoxb-\"")
	}

	return &ApplicationController{
		config: c,
		log:    l,
		client: socketmode.New(
			slack.New(
				c.WorkspaceToken,
				slack.OptionDebug(true),
				slack.OptionAppLevelToken(c.AppToken),
				slack.OptionLog(l),
			),
			socketmode.OptionDebug(true),
			socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
		),
	}
}

// ApplicationController is the handler for the slack events
type ApplicationController struct {
	config   *Config
	client   *socketmode.Client
	log      logs.Logger
	handlers map[EventType]EventHandler
}

// RegisterHandler for the incoming events
func (a *ApplicationController) RegisterHandler(ctx context.Context, eventType EventType, h EventHandler) {
	if a.handlers == nil {
		a.handlers = make(map[EventType]EventHandler)
	}
	a.handlers[eventType] = h
}

// Run the application starts listening to the socket and
// returns a cancel function to gracefully shutdown
func (a *ApplicationController) Run(ctx context.Context) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			a.log.Error(ctx, "Recovered from panic on application run.", logs.Error(r.(error)))
		}()

		for {
			select {
			case <-ctx.Done():
				a.log.Info(ctx, "Ending slack application")
				time.Sleep(time.Second)
				return
			case event := <-a.client.Events:
				switch event.Type {
				case socketmode.EventTypeConnecting:
					a.log.Info(ctx, "Connecting to Slack with Socket Mode...")
					break
				case socketmode.EventTypeConnectionError:
					a.log.Error(ctx, "Slack connection error trying to reconnect...")
					break
				case socketmode.EventTypeConnected:
					a.log.Error(ctx, "Connected to slack successfully...")
					break
				case socketmode.EventTypeEventsAPI:
					apiEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						// Ignore
						continue
					}

					a.log.Info(ctx, "Received an event from slack...")
					a.client.Ack(*event.Request)

					switch apiEvent.Type {
					case slackevents.CallbackEvent:
						eventType := EventType(apiEvent.InnerEvent.Type)
						h, ok := a.handlers[eventType]
						if !ok {
							a.log.Warn(ctx, fmt.Sprintf("Event of type '%s' was received but no handler has been found", event.Type))
							continue
						}

						err := h(context.Background(), apiEvent)
						if err != nil {
							a.log.Error(ctx, "An error has occurred executing the slack handler", logs.Error(err))
							return
						}
					default:
					}
					continue
				case socketmode.EventTypeInteractive:
				case socketmode.EventTypeSlashCommand:
				default:
					// Ignore
				}
				break
			default:
			}
		}
	}(ctx)

	err := a.client.Run()
	if err != nil {
		a.log.Error(ctx, "Couldn't initialize slack", logs.Error(err))
		return nil, err
	}

	return cancel, nil
}
