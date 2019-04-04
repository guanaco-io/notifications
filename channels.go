package main

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"strings"
)

type Channel interface {
	Send(event AlertEvent) error
}

type MailChannel struct {
	settings Smtp
	To       []string
	Template string
}

type SlackChannel struct {
	settings Slack
	Channel  string
}

type AlertEvent struct {
	NewAlertCount   int
	NewAlerts       []Alert
	AlreadyNotified int
}

func LoadChannels(config Config) (map[string]Channel, error) {

	var channels = make(map[string]Channel, len(config.Channels))

	for channelName, channel := range config.Channels {

		switch channel.Type {
		case "mail":
			tos, ok := channel.Config["to"]
			if !ok {
				return nil, errors.New(fmt.Sprintf("'to' property is required for channel '%v' of type 'mail' %v", channelName, channel.Type))
			}
			to := strings.Split(tos, ",")
			for i, t := range to {
				to[i] = strings.TrimSpace(t)
			}
			templateFilename, _ := channel.Config["templateFilename"]

			channels[channelName] = MailChannel{settings: config.ChannelSettings.Smtp, To: to, Template: templateFilename}

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

func (mail MailChannel) Send(event AlertEvent) error {
	log.Printf("Mailing %v alerts", event.NewAlertCount)

	body := render(getOrElse(mail.Template, "default_mail.gohtml"), event)

	m := gomail.NewMessage()
	m.SetHeader("From", mail.settings.From)
	m.SetHeader("To", mail.To...)
	m.SetHeader("Subject", subject(event.NewAlertCount))
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mail.settings.Server, mail.settings.Port(), mail.settings.User, mail.settings.Password)
	d.SSL = mail.settings.Ssl

	return d.DialAndSend(m)
}

func (slack SlackChannel) Send(event AlertEvent) error {
	log.Printf("Slacking %v alerts", event.NewAlertCount)

	return nil
}

func render(filename string, event AlertEvent) string {

	var result bytes.Buffer

	t := template.Must(template.New(filename).ParseFiles(fmt.Sprintf("config/%v", filename)))

	err := t.Execute(&result, event)
	if err != nil {
		panic(err)
	}

	return result.String()
}

func subject(count int) string {
	if count > 1 {
		return fmt.Sprintf("%v new alerts", count)
	}
	return fmt.Sprintf("%v new alert", count)
}

func getOrElse(attempt string, fallback string) string {
	if attempt == "" {
		return fallback
	}
	return attempt
}
