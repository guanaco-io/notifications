package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"testing"
	"time"
)

// An example showing how to unmarshal embedded
// structs from YAML.

var configData = `
alerta:
  endpoint: http://integratie.tjcc.be/alerta/api
  webui: http://integratie.tjcc.be/alerta
  reload_interval: 5
`


func TestReadFromYamlConfig(t *testing.T) {

	err := yaml.Unmarshal([]byte(configData), &Configuration)
	if err != nil {
		t.Fatalf("cannot unmarshal data: %v", err)
	}

	log.Printf("Parsed configuration: %v", Configuration)

	if (Configuration.Alerta.ReloadInterval != 5) {
		t.Fatalf("reload interval should equal 5")
	}

	log.Printf("interval is %v", Configuration.Alerta.ReloadInterval * time.Second)
}
