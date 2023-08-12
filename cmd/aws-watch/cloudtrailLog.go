package main

import "time"

type CloudtrailEvent struct {
	EventVersion                 string       `json:"eventVersion"`
	UserIdentity                 UserIdentity `json:"userIdentity"`
	EventTime                    time.Time    `json:"eventTime"`
	EventSource                  string       `json:"eventSource"`
	EventName                    string       `json:"eventName"`
	AwsRegion                    string       `json:"awsRegion"`
	SourceIPAddress              string       `json:"sourceIPAddress"`
	UserAgent                    string       `json:"userAgent"`
	RequestID                    string       `json:"requestID"`
	EventID                      string       `json:"eventID"`
	ReadOnly                     bool         `json:"readOnly"`
	EventType                    string       `json:"eventType"`
	ManagementEvent              bool         `json:"managementEvent"`
	RecipientAccountID           string       `json:"recipientAccountId"`
	EventCategory                string       `json:"eventCategory"`
	SessionCredentialFromConsole string       `json:"sessionCredentialFromConsole"`
	RequestParameters            interface{}  `json:"requestParameters"`
	ResponseElements             interface{}  `json:"responseElements"`
}
type SessionIssuer struct {
}
type WebIDFederationData struct {
}
type Attributes struct {
	CreationDate     time.Time `json:"creationDate"`
	MfaAuthenticated string    `json:"mfaAuthenticated"`
}
type SessionContext struct {
	SessionIssuer       SessionIssuer       `json:"sessionIssuer"`
	WebIDFederationData WebIDFederationData `json:"webIdFederationData"`
	Attributes          Attributes          `json:"attributes"`
}
type UserIdentity struct {
	Type           string         `json:"type"`
	PrincipalID    string         `json:"principalId"`
	Arn            string         `json:"arn"`
	AccountID      string         `json:"accountId"`
	AccessKeyID    string         `json:"accessKeyId"`
	UserName       string         `json:"userName"`
	SessionContext SessionContext `json:"sessionContext"`
}
