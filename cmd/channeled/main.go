package main

import (
	"github.com/gonzispina/channeled/kit/context"
	"github.com/gonzispina/channeled/kit/logs"
	"github.com/gonzispina/channeled/kit/slack"
	"github.com/slack-go/slack/slackevents"
	"os"
	"os/signal"
	"syscall"
)

var (
	// SIGTERM os.Signal
	SIGTERM os.Signal = syscall.SIGTERM
	// SIGTSTP os.Signal
	SIGTSTP os.Signal = syscall.SIGTSTP
	// SIGINT os.Signal
	SIGINT os.Signal = syscall.SIGINT
)

func main() {
	logs.InitDefault()
	log := logs.Log()
	ctx := context.Background()

	config := ReadConfig()
	controllers := NewControllersFactory(config, log)

	controllers.SlackAppController().RegisterHandler(ctx, slack.AppMention, func(ctx context.Context, event slackevents.EventsAPIEvent) error {
		log.Info(ctx, "Event received")
		return nil
	})

	cancelSlackApp, err := controllers.SlackAppController().Run(ctx)
	if err != nil {
		panic(err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, SIGINT)
	signal.Notify(sigchan, SIGTERM)
	signal.Notify(sigchan, SIGTSTP)

	<-sigchan

	cancelSlackApp()
}
