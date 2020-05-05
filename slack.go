package slack_command_proxy

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
)

type Message struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

type Payload struct {
	ChannelId     string `json:"channel_id"`
	ChannelName   string `json:"channel_name"`
	Command       string `json:"command"`
	ResponseUrl   string `json:"response_url"`
	TeamDomain    string `json:"team_domain"`
	TeamId        string `json:"team_id"`
	Text          string `json:"text"`
	Token         string `json:"token"`
	TriggerId     string `json:"trigger_id"`
	UserId        string `json:"user_id"`
	UserName      string `json:"user_name"`
	ParsedCommand string `json:"-"`
}

func (p *Payload) Bytes() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Payload.Bytes().json.Marshall(%v): %v", p, err)
	}

	return b
}

func newPayload(form url.Values) *Payload {
	b := formToJSONBytes(form)
	var p Payload

	err := json.Unmarshal(b, &p)
	if err != nil {
		log.Fatalf("newPayload().json.Unmarshall(%v, %v): %v", b, p, err)
	}

	p.ParsedCommand = strings.ReplaceAll(p.Command, "/", "")

	return &p
}

func formToJSONBytes(form url.Values) []byte {
	m := make(map[string]string)
	for key, values := range form {
		for _, value := range values {
			m[key] = value
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("newPayload().json.Marshall(%v): %v", m, err)
	}

	return b
}
