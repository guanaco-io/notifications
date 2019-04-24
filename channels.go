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
	"path"
	"strconv"
	"strings"
	"time"
)

type Channel interface {
	SendOpenAlerts(event OpenAlertsEvent, dryrun bool) error
	SendClosedAlerts(event ClosedAlertsEvent, dryrun bool) error
}

type MailChannel struct {
	Alerta   Alerta
	settings Smtp
	To       []string
	TemplateOpen string
	TemplateClosed string
}

type SlackChannel struct {
	Alerta   Alerta
	settings Slack
	Channel  string
}

type OpenAlertsEvent struct {
	NewAlertCount   int
	NewAlerts       []Alert
	AlreadyNotified int
}

type ClosedAlertsEvent struct {
	Alerts []Alert
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
			templateAlertsOpenedFilename, _ := channel.Config["template_open"]
			templateAlertsClosedFilename, _ := channel.Config["template_closed"]

			channels[channelName] = MailChannel{settings: config.ChannelSettings.Smtp, To: to, TemplateOpen: templateAlertsOpenedFilename, TemplateClosed:templateAlertsClosedFilename}

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

func (mail MailChannel) SendOpenAlerts(event OpenAlertsEvent, dryrun bool) error {

	mailTemplate := getOrElse(mail.TemplateOpen, "templates/open_alerts.gohtml")
	body := render(mailTemplate, event)

	return mail.send(event.Subject(), body, dryrun)
}

func (mail MailChannel) SendClosedAlerts(event ClosedAlertsEvent, dryrun bool) error {

	mailTemplate := getOrElse(mail.TemplateClosed, "templates/closed_alerts.gohtml")
	body := render(mailTemplate, event)

	return mail.send(event.Subject(), body, dryrun)
}

func (mail MailChannel) send(subject string, body string, dryrun bool) error {

	if dryrun {
		log.Print("-- DryRun is active: not really sending mail --")
		log.Printf("Generated mail from %v to %v [%v] \n%v", mail.settings.From, mail.To, subject, body)

		return nil
	} else {
		log.Printf("Sending smtp mail: %v", subject)

		m := gomail.NewMessage()
		m.SetHeader("From", mail.settings.From)
		m.SetHeader("To", mail.To...)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)

		d := gomail.NewDialer(mail.settings.Server, mail.settings.Port(), mail.settings.User, mail.settings.Password)
		d.SSL = mail.settings.Ssl

		return d.DialAndSend(m)
	}
}

func render(filename string, event interface{}) string {

	var result bytes.Buffer

	base := path.Base(filename)
	t := template.Must(template.New(base).ParseFiles(filename))

	err := t.Execute(&result, event)
	if err != nil {
		panic(err)
	}

	return result.String()
}

func (slackChannel SlackChannel) SendOpenAlerts(event OpenAlertsEvent, dryrun bool) error {

	msg := event.toWebhookMessage(slackChannel)

	return slackChannel.send(event.Subject(), msg, dryrun)
}

func (slackChannel SlackChannel) SendClosedAlerts(event ClosedAlertsEvent, dryrun bool) error {

	msg := event.toWebhookMessage(slackChannel)

	return slackChannel.send(event.Subject(), msg, dryrun)
}

func (slackChannel SlackChannel) send(subject string, body slack.WebhookMessage, dryrun bool) error {

	if dryrun {
		log.Print("-- DryRun is active: not really posting to slack --")

		if raw, err := json.Marshal(body); err != nil {
			log.Printf("Error marshalling slack message to json: %v", err)
			return err
		} else {
			log.Printf("Posting slack message:\n%v", string(raw))
			return nil
		}
	} else {
		log.Printf("Posting webhook msg to slack: %v", subject)
		return slack.PostWebhook(slackChannel.settings.WebhookUrl, &body)
	}
}

func (event OpenAlertsEvent) toWebhookMessage(slackChannel SlackChannel) slack.WebhookMessage {

	var attachments = make([]slack.Attachment, event.NewAlertCount)

	for index, alert := range event.NewAlerts {

		attachments[index] = slack.Attachment{
			Color: alert.Color(),
			//AuthorName: "Alerta Notifications",
			//AuthorLink: slackChannel.Alerta.Webui,

			Text: fmt.Sprintf("<%v|%v> - `%v` \n%v", alert.Url, alert.Resource, alert.Event, alert.Text),

			Footer: "Alerta Notifications",
			Ts:     json.Number(strconv.FormatInt(time.Now().Unix(), 10)),

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
		Text:        event.Subject(),
		Channel:     slackChannel.Channel,
		Attachments: attachments,
	}
	return msg
}

func (event ClosedAlertsEvent) toWebhookMessage(slackChannel SlackChannel) slack.WebhookMessage {

	var attachments = make([]slack.Attachment, len(event.Alerts))

	for index, alert := range event.Alerts {

		attachments[index] = slack.Attachment{
			Color: "#28a745",
			Text:  fmt.Sprintf("<%v|%v> - `%v`", alert.Url, alert.Resource, alert.Event),
		}
	}
	msg := slack.WebhookMessage{
		IconEmoji:   ":rocket:",
		Text:        event.Subject(),
		Channel:     slackChannel.Channel,
		Attachments: attachments,
	}
	return msg
}

func (event OpenAlertsEvent) Subject() string {

	if event.NewAlertCount > 1 {
		return fmt.Sprintf("%v new alerts", event.NewAlertCount)
	}
	return fmt.Sprintf("%v new alert", event.NewAlertCount)
}

func (event ClosedAlertsEvent) Subject() string {

	if len(event.Alerts) > 1 {
		return fmt.Sprintf("%v alerts were closed", len(event.Alerts))
	}
	return fmt.Sprintf("Alert %v was closed", event.Alerts[0].Resource)
}

func getOrElse(attempt string, fallback string) string {
	if attempt == "" {
		return fallback
	}
	return attempt
}
