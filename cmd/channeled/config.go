package main

import "github.com/gonzispina/channeled/kit/slack"

type Config struct {
	Slack slack.Config `json:"slack"`
}

// ReadConfig from the config file
func ReadConfig() *Config {
	return &Config{
		Slack: slack.Config{
			AppName:        "channeled",
			AppToken:       "xapp-1-A03BQ319FCJ-3395245544357-ac0c6c4fed82061809401a5f6523f7522a59c53800dfe9b9d7c7fc87583c3a76",
			WorkspaceToken: "xoxb-711669650370-3398240025730-3eL4JGEkQp7QAfk3v3FPhtLx",
		},
	}
}
