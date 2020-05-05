package slack_command_proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
)

var (
	config       Config
	projectId    = mustGetEnv("GCP_PROJECT")
	pubSubClient *pubsub.Client
)

func init() {
	// err is pre-declared to avoid shadowing client.
	var err error

	loadConfigFile()

	// client is initialized with context.Background() because it should
	// persist between function invocations.
	pubSubClient, err = pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		log.Fatalf("pubsub.NewClient(ctx, %s): %v", projectId, err)
	}
}

func mustGetEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

// SlackCommandProxy validates and publishes a Slack command to Cloud Pub/Sub.
func SlackCommandProxy(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading request body: %v", err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are accepted", http.StatusMethodNotAllowed)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		log.Fatalf("r.ParseForm(): %v", err)
	}

	var cmd *Command
	p := newPayload(r.Form)
	for _, c := range config.Commands {
		if c.TeamDomain == p.TeamDomain && c.Command == p.Command {
			cmd = &c
			break
		}
	}
	if cmd == nil {
		log.Fatalf("Error finding command '%s' in team '%s'", p.Command, p.TeamDomain)
	}

	// Reset r.Body as ParseForm depletes it by reading the io.ReadCloser.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	result, err := verifyRequest(r, cmd.SigningSecret)
	if err != nil {
		log.Fatalf("verifyRequest: %v", err)
	}
	if !result {
		log.Fatalf("Signatures did not match.")
	}

	if len(r.Form["text"]) == 0 {
		log.Fatalf("Empty text in form")
	}

	publish(w, r, p)

	var resp = &Message{
		ResponseType: "ephemeral",
		Text:         "I received your request!",
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		log.Fatalf("json.NewEncoder(w).Encode(%s): %v", resp, err)
	}
}
