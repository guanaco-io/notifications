package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var DryRun bool
var Configuration Config

type Alerta struct {

	Endpoint string `yaml:"endpoint"`
	Webui string `yaml:"webui"`
	ReloadInterval time.Duration `yaml:"reload_interval"`
}

type Config struct {

	Alerta Alerta `yaml:"alerta"`
}

func validate() error {

	if len(os.Args) < 2 {
		log.Printf("Usage: notifications [--dry-run] <config.yml>")
		log.Fatal("  <config.yml> parameter is missing!")
	}

	DryRun = contains(os.Args, "--dry-run")

	data, readFileError := ioutil.ReadFile(os.Args[1])
	if readFileError != nil {
		return readFileError
	}

	unmarshallError := yaml.Unmarshal(data, &Configuration)
	if unmarshallError != nil {
		return unmarshallError
	}

	return nil
}

func contains(slice []string, lookup string) bool {
	for _, element := range slice {
		if element == lookup {
			return true
		}
	}
	return false
}