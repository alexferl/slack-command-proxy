package slack_command_proxy

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	version                     = "v0"
	slackRequestTimestampHeader = "X-Slack-Request-Timestamp"
	slackSignatureHeader        = "X-Slack-Signature"
)

type oldTimeStampError struct {
	s string
}

func (e *oldTimeStampError) Error() string {
	return e.s
}

// verifyRequest verifies the request signature.
// See https://api.slack.com/docs/verifying-requests-from-slack.
func verifyRequest(r *http.Request, slackSigningSecret string) (bool, error) {
	timestamp := r.Header.Get(slackRequestTimestampHeader)
	slackSignature := r.Header.Get(slackSignatureHeader)

	if timestamp == "" {
		return false, fmt.Errorf("header '%s' is required", slackRequestTimestampHeader)
	}

	if slackSignature == "" {
		return false, fmt.Errorf("header '%s' is required", slackSignatureHeader)
	}

	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false, fmt.Errorf("strconv.ParseInt(%s): %v", timestamp, err)
	}

	if ageOk, age := checkTimestamp(t); !ageOk {
		return false, &oldTimeStampError{fmt.Sprintf("checkTimestamp(%v): %v %v", t, ageOk, age)}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, fmt.Errorf("ioutil.ReadAll(%v): %v", r.Body, err)
	}

	// Reset the body so other calls won't fail.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	baseString := fmt.Sprintf("%s:%s:%s", version, timestamp, body)

	signature := getSignature([]byte(baseString), []byte(slackSigningSecret))

	trimmed := strings.TrimPrefix(slackSignature, fmt.Sprintf("%s=", version))
	signatureInHeader, err := hex.DecodeString(trimmed)
	if err != nil {
		return false, fmt.Errorf("hex.DecodeString(%v): %v", trimmed, err)
	}

	return hmac.Equal(signature, signatureInHeader), nil
}

func getSignature(base []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(base)

	return h.Sum(nil)
}

// Arbitrarily trusting requests time stamped less than 5 minutes ago.
func checkTimestamp(timeStamp int64) (bool, time.Duration) {
	t := time.Since(time.Unix(timeStamp, 0))

	return t.Minutes() <= 5, t
}
