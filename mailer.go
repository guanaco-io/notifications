package main

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
)

func (channel MailChannel) Send(subject string, body string, dryrun bool) error {

	if dryrun {
		log.Print("-- DryRun is active: not really sending mail --")
		log.Printf("Generated mail from %v to %v [%v] \n%v", channel.settings.From, channel.To, subject, body)

		return nil
	} else {
		log.Printf("Sending smtp mail: %v", subject)

		if !channel.settings.Anonymous {

			m := gomail.NewMessage()
			m.SetHeader("From", channel.settings.From)
			m.SetHeader("To", channel.To...)
			m.SetHeader("Subject", subject)
			m.SetBody("text/html", body)

			d := gomail.NewDialer(channel.settings.Server, channel.settings.Port, channel.settings.User, channel.settings.Password)
			d.SSL = channel.settings.Ssl

			return d.DialAndSend(m)
		} else {
			// inspired by https://gadelkareem.com/2018/05/03/golang-send-mail-without-authentication-using-localhost-sendmail-or-postfix/
			to := make([]string, 0)
			for _, destination := range channel.To {
				to = append(to, (&mail.Address{"", destination}).String())
			}
			return SendAnonymous(
				fmt.Sprintf("%s:%d", channel.settings.Server, channel.settings.Port),
				(&mail.Address{channel.settings.FromName, channel.settings.From}).String(),
				subject,
				body,
				to,
			)
		}
	}
}

func SendAnonymous(addr, from, subject, body string, to []string) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(r.Replace(from)); err != nil {
		return err
	}
	for i := range to {
		to[i] = r.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
