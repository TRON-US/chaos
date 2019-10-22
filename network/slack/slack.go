package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachment struct {
	Color  string  `json:"color"`
	Fields []Field `json:"fields"`
}

type SlackRequestBody struct {
	Attachments []Attachment `json:"attachments"`
}

type Priority int

const (
	PRIORITY_0 Priority = iota
	PRIORITY_1
	PRIORITY_2
	PRIORITY_3
	PRIORITY_4
	PRIORITY_5
)

// @webhookUrl, the slack incoming webhookUrl to send the notification
// @slackRequest, the slack message body
// @retryTimeout, in seconds
// @priority, the priority of the slack message
// @priorityThreshold, the priority threshold determining whether to send this slack message based on the priority level
func SendSlackNotification(webhookUrl string, slackRequest SlackRequestBody, retryTimeout int, priority Priority, priorityThreshold Priority) error {
	if priority > priorityThreshold {
		return nil
	}

	// marshal slack request and create slack webhook request
	slackBody, _ := json.Marshal(slackRequest)
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}

	// make slack webhook request with retry
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Duration(retryTimeout) * time.Second
	var resp *http.Response
	sendNotification := func() error {
		resp, err = client.Do(req)
		return err
	}
	err = backoff.Retry(sendNotification, b)
	if err != nil {
		return err
	}

	// read response of slack webhook request
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if buf.String() != "ok" {
		return errors.New("non-ok response returned from Slack")
	}
	return nil
}
