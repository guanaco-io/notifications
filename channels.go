package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"strings"
)

type Channel interface {
	Send(event AlertEvent, dryrun bool) error
}

type MailChannel struct {
	Alerta   Alerta
	settings Smtp
	To       []string
	Template string
}

type SlackChannel struct {
	Alerta   Alerta
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

func (mail MailChannel) Send(event AlertEvent, dryrun bool) error {
	log.Printf("Mailing %v alerts", event.NewAlertCount)

	body := render(getOrElse(mail.Template, "default_mail.gohtml"), event)

	if dryrun {
		log.Print("-- DryRun is active: not really sending mail --")
		log.Printf("Generated mail from %v to %v [%v] \n%v", mail.settings.From, mail.To, subject(event.NewAlertCount), body)

		return nil
	} else {

		m := gomail.NewMessage()
		m.SetHeader("From", mail.settings.From)
		m.SetHeader("To", mail.To...)
		m.SetHeader("Subject", subject(event.NewAlertCount))
		m.SetBody("text/html", body)

		d := gomail.NewDialer(mail.settings.Server, mail.settings.Port(), mail.settings.User, mail.settings.Password)
		d.SSL = mail.settings.Ssl

		return d.DialAndSend(m)
	}
}

func render(filename string, event AlertEvent) string {

	var result bytes.Buffer

	t := template.Must(template.New(filename).ParseFiles(fmt.Sprintf("templates/%v", filename)))

	err := t.Execute(&result, event)
	if err != nil {
		panic(err)
	}

	return result.String()
}

func (slackChannel SlackChannel) Send(event AlertEvent, dryrun bool) error {
	log.Printf("Slacking %v alerts", event.NewAlertCount)

	msg := toWebhookMessage(event, slackChannel)

	if dryrun {
		log.Print("-- DryRun is active: not really posting to slack --")

		if raw, err := json.Marshal(msg); err != nil {
			log.Printf("Error marshalling slack message to json: %v", err)
			return err
		} else {
			log.Printf("Posting slack message:\n%v", string(raw))
			return nil
		}
	} else {
		return slack.PostWebhook(slackChannel.settings.WebhookUrl, &msg)
	}
}

func toWebhookMessage(event AlertEvent, slackChannel SlackChannel) slack.WebhookMessage {

	var attachments = make([]slack.Attachment, event.NewAlertCount)

	for index, alert := range event.NewAlerts {

		attachments[index] = slack.Attachment{
			Color:      color(alert.Severity),
			AuthorName: "Alerta Notifications",
			AuthorLink: slackChannel.Alerta.Webui,

			Text: fmt.Sprintf("*<%v|%v>* - `%v`\n%v", alert.Url, alert.Event, alert.Resource, alert.Text),

			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "Severity",
					Value: alert.Severity,
					Short: true,
				},
				slack.AttachmentField{
					Title: "Environment",
					Value: alert.Environment,
					Short: true,
				},
			},
		}
	}
	msg := slack.WebhookMessage{
		IconEmoji:   ":rocket:",
		Text:        fmt.Sprintf(subject(event.NewAlertCount)),
		Channel:     slackChannel.Channel,
		Attachments: attachments,
	}
	return msg
}

func subject(count int) string {
	if count > 1 {
		return fmt.Sprintf("%v new alerts", count)
	}
	return fmt.Sprintf("%v new alert", count)
}

func color(severity string) string {
	switch severity {
	case "warning":
		return "danger"
	case "minor":
		return "warning"
	default:
		return "danger"
	}
}

func getOrElse(attempt string, fallback string) string {
	if attempt == "" {
		return fallback
	}
	return attempt
}
