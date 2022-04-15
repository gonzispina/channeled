package slack

import (
	"github.com/gonzispina/channeled/kit/context"
	"github.com/slack-go/slack/slackevents"
)

type EventType string

const (
	// AppMention the application was mentioned by tag
	AppMention EventType = "app_mention"
	// DirectMessageReceived a message was posted in a direct message channel
	DirectMessageReceived EventType = "message.im"
	// PrivateGroupMessageReceived a message was posted to a private channel
	PrivateGroupMessageReceived EventType = "message.groups"
	// PublicGroupMessageReceived a message was posted to a public channel
	PublicGroupMessageReceived EventType = "message.channels"
)

// EventHandler function to execute
type EventHandler func(ctx context.Context, event slackevents.EventsAPIEvent) error
