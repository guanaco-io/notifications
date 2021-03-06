package main

import (
	"log"
	"testing"
	"time"
)

func TestReadFromYamlConfig(t *testing.T) {

	Configuration, configError := Load("config/config.yml")

	if configError != nil {
		t.Fatalf("cannot unmarshal data: %v", configError)
	}
	log.Printf("Parsed configuration: %v\n", Configuration)

	if Configuration.DryRun != true {
		t.Fatalf("unexpected dry_run configuration value")
	}
	log.Printf("dryrun is %v", Configuration.DryRun)

	if Configuration.Alerta.ReloadInterval != 60 {
		t.Fatalf("unexpected reload interval")
	}
	log.Printf("interval is %v", Configuration.Alerta.ReloadInterval*time.Second)

	log.Printf("channels: %v", Configuration.Channels)
	if len(Configuration.Channels) != 3 {
		t.Fatalf("expected 3 channels")
	}
	for name, channel := range Configuration.Channels {
		log.Printf("channel[%v]: type '%v', config %v", name, channel.Type, channel.Config)

		if channel.Type == "mail" {
			log.Printf("channel[%v].template_open: '%v'", name, channel.Config["template_open"])
			log.Printf("channel[%v].template_closed: '%v'", name, channel.Config["template_closed"])
		}
	}

	log.Printf("rules: %v", Configuration.Rules)
	if len(Configuration.Rules) != 2 {
		t.Fatalf("expected 2 rules")
	}

	for rulename, rule := range Configuration.Rules {
		log.Printf("rule[%v]: filter '%v', channels %v", rulename, rule.Filter, rule.Channels)
	}
}
