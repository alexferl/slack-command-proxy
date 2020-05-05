package slack_command_proxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Commands []Command
}

type Command struct {
	Command       string `json:"command"`
	SigningSecret string `json:"signing_secret"`
	TeamDomain    string `json:"team_domain"`
}

func loadConfigFile() {
	file, err := ioutil.ReadFile("./serverless_function_source_code/config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}
}
