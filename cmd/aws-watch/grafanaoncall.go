package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

type GrafanaOncallMessage struct {
	UserIdentity      UserIdentity `json:"userIdentity"`
	EventTime         time.Time    `json:"eventTime"`
	EventSource       string       `json:"eventSource"`
	EventName         string       `json:"eventName"`
	AwsRegion         string       `json:"awsRegion"`
	SourceIPAddress   string       `json:"sourceIPAddress"`
	UserAgent         string       `json:"userAgent"`
	EventID           string       `json:"eventID"`
	EventCategory     string       `json:"eventCategory"`
	RequestParameters interface{}  `json:"requestParameters"`
	ResponseElements  interface{}  `json:"responseElements"`
}

type GrafanaOncallEKSMessage struct {
	RequestURI     string            `json:"requestURI"`
	Verb           string            `json:"verb"`
	User           EKSUser           `json:"user"`
	SourceIPs      []string          `json:"sourceIPs"`
	UserAgent      string            `json:"userAgent"`
	ObjectRef      EKSObjectRef      `json:"objectRef"`
	ResponseStatus EKSResponseStatus `json:"responseStatus"`
	StageTimestamp time.Time         `json:"stageTimestamp"`
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

func (g *GrafanaOncallClient) buildMessage(e Event) (string, error) {
	switch e.EventSource {
	case "EKS":
		return g.buildEKSMessage(*e.EKS)
	default:
		return g.buildCloudtrailMessage(*e.Cloudtrail)
	}
}

func (g *GrafanaOncallClient) buildCloudtrailMessage(e CloudtrailEvent) (string, error) {
	msg := &GrafanaOncallMessage{
		UserIdentity:      e.UserIdentity,
		EventTime:         e.EventTime,
		EventSource:       e.EventSource,
		EventName:         e.EventName,
		AwsRegion:         e.AwsRegion,
		SourceIPAddress:   e.SourceIPAddress,
		UserAgent:         e.UserAgent,
		EventID:           e.EventID,
		EventCategory:     e.EventCategory,
		RequestParameters: e.RequestParameters,
		ResponseElements:  e.ResponseElements,
	}

	jsonStr, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	return string(jsonStr), err
}

func (g *GrafanaOncallClient) buildEKSMessage(e EKSEvent) (string, error) {
	msg := &GrafanaOncallEKSMessage{
		RequestURI:     e.RequestURI,
		Verb:           e.Verb,
		User:           e.User,
		SourceIPs:      e.SourceIPs,
		UserAgent:      e.UserAgent,
		ObjectRef:      e.ObjectRef,
		ResponseStatus: e.ResponseStatus,
		StageTimestamp: e.StageTimestamp,
	}

	jsonStr, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	return string(jsonStr), err
}
