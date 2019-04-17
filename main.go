package main

import (
	"log"
	"os"
	"time"
)

const (
	maxDuration time.Duration = 1<<63 - 1
)

func main() {

	log.Printf("Starting Guanaco notifications app")
	if len(os.Args) < 2 {
		log.Printf("Usage: notifications <config.yml>")
		log.Fatal("  <config.yml> parameter is missing!")
	}

	config, initError := Load(os.Args[1])
	logFatal("Error initializing program", initError)
	log.Printf("Configuration loaded successfully")

	channels, channelsError := LoadChannels(config)
	logFatal("Error loading channels configuration", channelsError)
	log.Printf("%v Channels loaded successfully", len(channels))

	client := AlertaClient{config: config.Alerta}
	ticker := time.NewTicker(config.Alerta.ReloadInterval * time.Second)

	log.Printf("Waiting for %v before fetching alerts", config.Alerta.ReloadInterval*time.Second)
	go func() {
		for t := range ticker.C {

			for ruleName, rule := range config.Rules {

				log.Printf("Evaluating rule %v (%v)", ruleName, t)

				if alerts := client.searchAlerts(rule); alerts != nil && len(alerts) > 0 {

					alreadyNotified, notNotified := Partition(alerts, ruleName, IsNotified)

					if nrOfAlerts := len(notNotified); nrOfAlerts > 0 {

						for _, ruleChannel := range rule.Channels {
							log.Printf("Sending %v alert(s) to channel %v of rule %v", nrOfAlerts, ruleChannel, ruleName)

							channel, ok := channels[ruleChannel]
							if !ok {
								log.Fatalf("Unable to find channel '%v' of rule '%v' in channel config", ruleChannel, ruleName)
							}

							sendError := channel.Send(AlertEvent{NewAlertCount: nrOfAlerts, NewAlerts: notNotified, AlreadyNotified: len(alreadyNotified)}, config.DryRun)
							if sendError != nil {
								log.Printf("Error sending alert event to channel '%v' of rule '%v': %v", ruleChannel, ruleName, sendError)
							}
						}

						for _, alert := range notNotified {
							alert.Notified(ruleName)
							updateError := client.updateAttributes(alert, config.DryRun)
							if updateError != nil {
								log.Printf("Error updating alert attributes for alert '%v' and rule '%v': %v", alert, ruleName, updateError)
							}
						}
					}
					log.Printf("%v alerts were already notified for rule %v", len(alreadyNotified), ruleName)

				} else {
					log.Printf("No Alerts found for rule %v", ruleName)
				}

			}
		}
	}()

	time.Sleep(maxDuration) // TODO find a better way to block here
	ticker.Stop()
	log.Println("Ticker stopped, exiting program")

	os.Exit(0)
}

func Partition(all []Alert, ruleId string, predicate func(Alert, string) bool) ([]Alert, []Alert) {
	success := make([]Alert, 0)
	failure := make([]Alert, 0)
	for _, alert := range all {
		if predicate(alert, ruleId) {
			success = append(success, alert)
		} else {
			failure = append(failure, alert)
		}
	}
	return success, failure
}

// Log error message and exit program
func logFatal(msg string, err error) {
	if err != nil {
		log.Fatalf("%v: %v", msg, err)
	}
}
