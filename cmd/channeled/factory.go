package main

import (
	"github.com/gonzispina/channeled/kit/logs"
	"github.com/gonzispina/channeled/kit/slack"
)

// NewControllersFactory constructor
func NewControllersFactory(c *Config, l logs.Logger) *ControllersFactory {
	if c == nil {
		panic("config must be initialized")
	}

	if l == nil {
		panic("logger must be initialized")
	}

	return &ControllersFactory{config: c, log: l}
}

type ControllersFactory struct {
	config *Config
	log    logs.Logger

	slackAppController *slack.ApplicationController
}

func (f *ControllersFactory) SlackAppController() *slack.ApplicationController {
	if f.slackAppController == nil {
		f.slackAppController = slack.NewApplication(
			&f.config.Slack,
			f.log,
		)
	}
	return f.slackAppController
}
