package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	cfgURL         = "internal.webhooks.ente.url"
	cfgToken       = "internal.webhooks.ente.token"
	headerToken    = "X-Ente-Token"
	defaultTimeout = 5 * time.Second
)

type payload struct {
	Event string `json:"event"`
	Email string `json:"email"`
}

// Send posts a webhook payload to the configured Ente webhook URL.
func Send(event string, email string) {
	url := viper.GetString(cfgURL)
	if url == "" {
		return
	}

	body, err := json.Marshal(payload{
		Event: event,
		Email: email,
	})
	if err != nil {
		log.WithError(err).Warn("Send webhook: failed to marshal payload")
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.WithError(err).Warn("Send webhook: failed to create request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	token := viper.GetString(cfgToken)
	if token != "" {
		req.Header.Set(headerToken, token)
	}

	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).WithField("url", url).Warn("Send webhook: request failed")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.WithFields(log.Fields{
			"url":    url,
			"status": resp.StatusCode,
		}).Warn("Send webhook: non-2xx response")
	}
}
