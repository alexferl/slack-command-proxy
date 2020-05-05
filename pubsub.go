package slack_command_proxy

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
)

func publish(w http.ResponseWriter, r *http.Request, p *Payload) {
	topicName := fmt.Sprintf("slack-command-proxy-%s-%s", p.TeamDomain, p.ParsedCommand)

	msg := &pubsub.Message{
		Data: p.Bytes(),
	}

	if _, err := pubSubClient.Topic(topicName).Publish(r.Context(), msg).Get(r.Context()); err != nil {
		http.Error(w, fmt.Sprintf("Error publishing message: %v", err), http.StatusInternalServerError)
		return
	}
}
