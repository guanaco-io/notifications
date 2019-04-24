package main

import (
	"log"
	"time"
)

type RuleHandler struct {
	alerta AlertaClient

	ruleName string
	rule     Rule

	channels map[string]Channel

	openAlerts []Alert

	dryRun bool
}

func (handler *RuleHandler) handle(time time.Time) {

	log.Printf("Evaluating rule %v (%v)", handler.ruleName, time)

	if openAlerts := handler.alerta.searchAlerts(handler.rule); openAlerts != nil && len(openAlerts) > 0 {

		alreadyNotified, notNotified := Partition(openAlerts, handler.ruleName, IsNotified)

		if nrOfAlerts := len(notNotified); nrOfAlerts > 0 {

			for _, ruleChannel := range handler.rule.Channels {
				log.Printf("Sending %v alert(s) to channel %v of rule %v", nrOfAlerts, ruleChannel, handler.ruleName)

				channel, ok := handler.channels[ruleChannel]
				if !ok {
					log.Fatalf("Unable to find channel '%v' of rule '%v' in channel config", ruleChannel, handler.ruleName)
				}

				sendError := channel.SendOpenAlerts(OpenAlertsEvent{NewAlertCount: nrOfAlerts, NewAlerts: notNotified, AlreadyNotified: len(alreadyNotified)}, handler.dryRun)
				if sendError != nil {
					log.Printf("Error sending alert event to channel '%v' of rule '%v': %v", ruleChannel, handler.ruleName, sendError)
				}
			}

			for _, alert := range notNotified {
				alert.Notified(handler.ruleName)
				updateError := handler.alerta.updateAttributes(alert, handler.dryRun)
				if updateError != nil {
					log.Printf("Error updating alert attributes for alert '%v' and rule '%v': %v", alert, handler.ruleName, updateError)
				}
			}
		}
		log.Printf("%v alerts were already notified for rule %v", len(alreadyNotified), handler.ruleName)

		if closedAlerts := handler.getClosedAlerts(openAlerts); closedAlerts != nil && len(closedAlerts) > 0 {

			log.Printf("%v alerts were closed for rule %v", len(closedAlerts), handler.ruleName)

			for _, ruleChannel := range handler.rule.Channels {
				log.Printf("Sending %v closed alert(s) to channel %v of rule %v", len(closedAlerts), ruleChannel, handler.ruleName)

				channel, ok := handler.channels[ruleChannel]
				if !ok {
					log.Fatalf("Unable to find channel '%v' of rule '%v' in channel config", ruleChannel, handler.ruleName)
				}

				sendError := channel.SendClosedAlerts(ClosedAlertsEvent{Alerts: closedAlerts}, handler.dryRun)
				if sendError != nil {
					log.Printf("Error sending closed alerts event to channel '%v' of rule '%v': %v", ruleChannel, handler.ruleName, sendError)
				}
			}
		} else {
			log.Printf("%v alerts were closed for rule %v", len(closedAlerts), handler.ruleName)
		}

		handler.openAlerts = openAlerts
		log.Printf("tracking %v open alerts for rule %v", len(handler.openAlerts), handler.ruleName)

	} else {
		log.Printf("No Alerts found for rule %v", handler.ruleName)
	}
}

func (handler *RuleHandler) getClosedAlerts(currentOpenAlerts []Alert) []Alert {

	closedAlerts := make([]Alert, 0)

	for _, previouslyOpenAlert := range handler.openAlerts {

		if !Contains(previouslyOpenAlert, currentOpenAlerts) {
			closedAlerts = append(closedAlerts, previouslyOpenAlert)
		}
	}

	return closedAlerts
}

func Contains(alert Alert, alerts []Alert) bool {
	for _, candidate := range alerts {
		if candidate.Id == alert.Id {
			return true
		}
	}
	return false
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
