package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestLoadChannelsConfig(t *testing.T) {

	Configuration, configError := Load("config/config.yml")
	if configError != nil {
		t.Fatalf("cannot unmarshal data: %v", configError)
	}

	Channels, channelLoadError := LoadChannels(Configuration)
	if channelLoadError != nil {
		t.Fatalf("cannot load channels from config: %v", channelLoadError)
	}

	if len(Channels) != 3 {
		t.Fatalf("expected 3 channels")
	}
	for name, channel := range Channels {
		log.Printf("channel[%v]: %T %v", name, channel, channel)
	}
}

func TestMailTemplate(t *testing.T) {

	mockAlertEvent := AlertEvent{AlreadyNotified: 20, NewAlertCount: 5, NewAlerts: readAlerts(t)}

	log.Print(render("default_mail.gohtml", mockAlertEvent))
}

func TestSlackMarshalling(t *testing.T) {

	mockAlertEvent := AlertEvent{AlreadyNotified: 20, NewAlertCount: 5, NewAlerts: readAlerts(t)}
	mockChannel := SlackChannel{Channel: "test", Alerta: Alerta{Webui: "http://localhost:8282/alerta"}}

	msg := toWebhookMessage(mockAlertEvent, mockChannel)

	raw, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("cannot marshall webhook message: %v", err)
	}

	log.Print(string(raw))
}

func readAlerts(t *testing.T) []Alert {
	var alertsResponse AlertsResponse

	data, readFileError := ioutil.ReadFile("test/alerts.json")
	if readFileError != nil {
		t.Fatalf("cannot read test file test/alerts.json: %v", readFileError)
	}

	err := json.Unmarshal(data, &alertsResponse)
	if err != nil {
		t.Fatalf("cannot unmarshal data: %v", err)
	}

	for index, alert := range alertsResponse.Alerts {
		alertsResponse.Alerts[index].Url = fmt.Sprintf("%v/#/alert/%v", "http://localhost:8282", alert.Id)
	}

	return alertsResponse.Alerts
}
