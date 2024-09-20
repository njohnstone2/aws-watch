package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type GrafanaOncallClient struct {
	client *http.Client
}

func NewGrafanaOncallClient() *GrafanaOncallClient {
	return &GrafanaOncallClient{
		client: &http.Client{},
	}
}

type Alert struct {
	UID          string `json:"alert_uid"`
	Title        string `json:"title"`
	ImageUrl     string `json:"image_url"`
	State        string `json:"state"`
	UpstreamLink string `json:"link_to_upstream_details"`
	Message      string `json:"message"`
}

func (g *GrafanaOncallClient) CreateAlert(url string, a *Alert) error {
	jsonStr, err := json.Marshal(a)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to create alert. StatusCode: %d", resp.StatusCode)
	}

	return nil
}

func (g *GrafanaOncallClient) buildMessage(e Event) string {
	switch e.EventSource {
	case "EKS":
		return g.buildEKSMessage(*e.EKS)
	default:
		return g.buildCloudtrailMessage(*e.Cloudtrail)
	}
}

func (g *GrafanaOncallClient) buildCloudtrailMessage(e CloudtrailEvent) string {
	requestObj := prettyPrintJson(e.RequestParameters)

	msg := fmt.Sprintf("AWS configuration change detected in account `%s` \n- Service: `%s` \n- Username: %s \n- Action: %s \n- Region: %s \n- UserAgent: %s \n- Time: %s \n- Resource: ```%s```",
		e.UserIdentity.AccountID,
		e.EventSource,
		e.UserIdentity.UserName,
		e.EventName,
		e.AwsRegion,
		e.UserAgent,
		e.EventTime.Format("2006-01-02 15:04:05"),
		requestObj,
	)

	return msg
}

func (g *GrafanaOncallClient) buildEKSMessage(e EKSEvent) string {
	msg := fmt.Sprintf("EKS %s change detected in cluster `%s` \n- Username: %s \n- Name: %s \n- Namespace: %s \n- URI: `%s` \n- Time: %s",
		e.ObjectRef.Resource,
		e.ClusterName,
		e.User.Username,
		e.ObjectRef.Name,
		e.ObjectRef.Namespace,
		e.RequestURI,
		e.RequestReceivedTimestamp.Format("2006-01-02 15:04:05"),
	)

	if e.RequestObject.Data != nil {
		requestObj := prettyPrintJson(e.RequestObject.Data)

		msg += fmt.Sprintf("\nRequest:\n```%s```", requestObj)
	}

	if e.ResponseObject.Data != nil {
		responseObj := prettyPrintJson(e.ResponseObject.Data)

		msg += fmt.Sprintf("\nResponse:\n```%s```", responseObj)
	}

	return msg
}

func prettyPrintJson(d interface{}) string {
	b, err := json.MarshalIndent(d, "", "  ")
	log.WithError(err).Error("prettyPrintJson() failed to parse object")

	return string(b)
}
