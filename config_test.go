package main

import (
	"log"
	"testing"
	"time"
)

func TestReadFromYamlConfig(t *testing.T) {

	Configuration, configError := Load("config/config.yaml")

	if configError != nil {
		t.Fatalf("cannot unmarshal data: %v", configError)
	}
	log.Println(Configuration)

	log.Printf("Parsed configuration: %v\n", Configuration)

	if Configuration.Alerta.ReloadInterval != 60 {
		t.Fatalf("unexpected reload interval")
	}
	log.Printf("interval is %v\n", Configuration.Alerta.ReloadInterval*time.Second)

	log.Printf("channels: %v\n", Configuration.Channels)
	if len(Configuration.Channels) != 3 {
		t.Fatalf("expected 3 channels")
	}
	for name, channel := range Configuration.Channels {
		log.Printf("channel[%v]: type '%v', config %v", name, channel.Type, channel.Config)
	}

	log.Printf("rules: %v\n", Configuration.Rules)
	if len(Configuration.Rules) != 2 {
		t.Fatalf("expected 2 rules")
	}

	for rulename, rule := range Configuration.Rules {
		log.Printf("rule[%v]: filter '%v', channels %v", rulename, rule.Filter, rule.Channels)
	}
}
