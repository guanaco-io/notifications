package main

import (
	"errors"
	"fmt"
	"log"
)

type Channel interface {

	Send(alerts []Alert)
}

type MailChannel struct {

	settings Smtp
	To string
	Template string
}

type SlackChannel struct {

	settings Slack
	Channel string
}

func LoadChannels(config Config) (map[string] Channel, error) {

	var channels = make(map[string]Channel, len(config.Channels))

	for channelName, channel := range config.Channels {

		switch channel.Type {
			case "mail":
				to, ok := channel.Config["to"]
				if !ok {
					return nil, errors.New(fmt.Sprintf("'to' property is required for channel '%v' of type 'mail' %v", channelName, channel.Type))
				}
				template, _ := channel.Config["template"]

				channels[channelName] = MailChannel{settings: config.ChannelSettings.Smtp, To: to, Template: template}

			case "slack":
				slackChannel, ok := channel.Config["slack_channel"]
				if !ok {
					return nil, errors.New(fmt.Sprintf("'channel' property is required for channel '%v' of type 'slack' %v", channelName, channel.Type))
				}
				channels[channelName] = SlackChannel{settings: config.ChannelSettings.Slack, Channel: slackChannel}

			default:
				return nil, errors.New(fmt.Sprintf("Unknown channel type %v: valid types are %v", channel.Type, "mail, slack"))
		}
	}
	return channels, nil
}

func (mail MailChannel) Send(alerts []Alert) {
	log.Printf("Mailing %v alerts", len(alerts))
}

func (slack SlackChannel) Send(alerts []Alert) {
	log.Printf("Slacking %v alerts", len(alerts))
}