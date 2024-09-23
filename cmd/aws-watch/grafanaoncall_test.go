package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestCreateAlert(t *testing.T) {
	t.Run("Create Alert", func(t *testing.T) {

		oncallClient := &GrafanaOncallClient{
			client: &MockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					if req.URL.Path != "/integrations/v1/formatted_webhook/6FJiWRoaASVbsMYrUPy5XQq0V/" {
						t.Errorf("Expected to request '/integrations/v1/formatted_webhook/6FJiWRoaASVbsMYrUPy5XQq0V/', got: %s", req.URL.Path)
					}

					responseBody := io.NopCloser(bytes.NewReader([]byte(`"Ok."`)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
			},
		}

		WebhookUrl := "http://localhost:8080/integrations/v1/formatted_webhook/6FJiWRoaASVbsMYrUPy5XQq0V/"
		alert := &Alert{
			UID:     uuid.New().String(),
			Title:   "Audit Event (EKS)",
			State:   "alerting",
			Message: "EKS configmaps change detected in cluster `example-eks` \n- Username: aad:test.user@example.com \n- Name: example-configmap \n- Namespace: default \n- URI: `/api/v1/namespaces/default/configmaps/example-configmap?fieldManager=kubectl-edit&fieldValidation=Strict` \n- Time: 2024-01-02 03:04:05",
		}

		err := oncallClient.CreateAlert(WebhookUrl, alert)
		assert.NoError(t, err)
	})
}

func TestBuildMessage(t *testing.T) {
	t.Run("Build Cloudtrail Message", func(t *testing.T) {
		cloudtrailEvent := &Event{
			EventSource: "ec2.amazonaws.com/",
			Cloudtrail:  &cloudtrailEvent,
		}

		expected := "AWS configuration change detected in account `123456789012` \n- Service: `ec2.amazonaws.com` \n- Username: example \n- Action: CreateSecurityGroup \n- Region: us-east-1 \n- UserAgent: AWS Internal \n- Time: 0001-01-01 00:00:00 \n- Resource: ```\"{\\\"groupName\\\":\\\"temp\\\",\\\"groupDescription\\\":\\\"example audit event\\\",\\\"vpcId\\\":\\\"vpc-a11a1111\\\"}\"```"

		client := NewGrafanaOncallClient()

		msg := client.buildMessage(*cloudtrailEvent)
		assert.Equal(t, expected, msg)
	})

	t.Run("Build EKS Message", func(t *testing.T) {
		eksEvent := &Event{
			EventSource: "EKS",
			EKS:         &eksEvent,
		}

		expected := "EKS configmaps change detected in cluster `example-eks` \n- Username: example-user \n- Name: example \n- Namespace: default \n- URI: `/api/v1/namespaces/default/configmaps/example?fieldManager=kubectl-edit&fieldValidation=Strict` \n- Time: 0001-01-01 00:00:00\nRequest:\n```{\n  \"first\": \"abc\",\n  \"second\": \"def\"\n}```\nResponse:\n```{\n  \"first\": \"abc\",\n  \"second\": \"def\",\n  \"third\": \"123\"\n}```"

		client := NewGrafanaOncallClient()

		msg := client.buildMessage(*eksEvent)
		assert.Equal(t, expected, msg)
	})
}
