package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const notification_attribute_format = "notifications %s"

type AlertaClient struct {
	config Alerta
}

type Alert struct {
	Id          string            `json:"id"`
	Environment string            `json:"environment"`
	Text        string            `json:"text"`
	Resource    string            `json:"resource"`
	Severity    string            `json:"severity"`
	Event       string            `json:"event"`
	Href        string            `json:"href"`
	Attributes  map[string]string `json:"attributes"`

	Url string
}

type AlertsResponse struct {
	Alerts       []Alert        `json:"alerts"`
	StatusCounts map[string]int `json:"statusCounts"`
}

func (alert *Alert) AlreadyNotified(ruleId string) bool {
	_, ok := alert.Attributes[fmt.Sprintf(notification_attribute_format, ruleId)]
	return ok
}

func (alert *Alert) Notified(ruleId string) {
	alert.Attributes[fmt.Sprintf(notification_attribute_format, ruleId)] = time.Now().UTC().String()
}

func (alert *Alert) Color() string {
	switch alert.Severity {
	case "warning":
		return "#17a2b8"
	case "minor":
		return "#ffc107"
	case "major":
		return "#dc3545"
	case "critical":
		return "#dc3545"
	default:
		return "#343a40"
	}
}

func IsNotified(alert Alert, ruleId string) bool {
	return alert.AlreadyNotified(ruleId)
}

func (client *AlertaClient) searchAlerts(rule Rule) []Alert {
	var alertsResponse = AlertsResponse{}

	url := fmt.Sprintf("%v/alerts?%v", client.config.Endpoint, rule.Filter)
	resp, err := client.performRequest("GET", url, nil)

	if err != nil {
		log.Printf("Error fetching alerts: %v", err)
		return alertsResponse.Alerts
	}

	log.Printf("< %s", resp.Status)

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&alertsResponse); err != nil {
		log.Fatalf("Error parsing alerts response: %v", err)
		return alertsResponse.Alerts
	}

	for index, alert := range alertsResponse.Alerts {
		alertsResponse.Alerts[index].Url = fmt.Sprintf("%v/#/alert/%v", client.config.Webui, alert.Id)
	}

	closeError := resp.Body.Close()
	if closeError != nil {
		log.Fatalf("Error closing response body: %v", closeError)
	}

	return alertsResponse.Alerts
}

// http://docs.alerta.io/en/latest/api/reference.html#update-alert-attributes
func (client *AlertaClient) updateAttributes(alert Alert, dryrun bool) error {

	url := fmt.Sprintf("%v/alert/%v/attributes", client.config.Endpoint, alert.Id)

	var body = make(map[string]interface{})
	body["attributes"] = alert.Attributes

	jsn, marshallError := json.Marshal(body)
	if marshallError != nil {
		return marshallError
	}

	if dryrun {
		log.Print("-- DryRun is active: not really updating alert attributes --")
		log.Printf("Generated attribute update request for Alerta API: [%v] \n%v", url, string(jsn))

		return nil
	} else {
		log.Printf("Posting attribute update to Alerta")

		response, err := client.performRequest("PUT", url, jsn)
		if err == nil {
			log.Printf("< %v", response.Status)
		}
		return err
	}
}

func (alerta *AlertaClient) performRequest(method string, url string, body []byte) (resp *http.Response, err error) {
	log.Printf("> [%s] %s", method, url)
	if body != nil {
		log.Printf("> %s", body)
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if alerta.config.ApiToken != "" {
		req.Header.Set("X-API-Key", alerta.config.ApiToken)
	}
	client := &http.Client{}
	return client.Do(req)
}
