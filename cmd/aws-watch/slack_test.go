package main

import (
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

var cloudtrailEvent = CloudtrailEvent{
	EventVersion: "1.08",
	UserIdentity: UserIdentity{
		Type:        "IAMUser",
		PrincipalID: "AIDA1AAAAAAAAA1AA1AAA",
		Arn:         "arn:aws:iam::123456789012:user/example",
		AccountID:   "123456789012",
		AccessKeyID: "ASIA1AAAAAAAAAAAA1AA",
		UserName:    "example",
	},
	EventTime:         time.Time{},
	EventSource:       "ec2.amazonaws.com",
	EventName:         "CreateSecurityGroup",
	AwsRegion:         "us-east-1",
	SourceIPAddress:   "1.2.3.4",
	UserAgent:         "AWS Internal",
	RequestParameters: "{\"groupName\":\"temp\",\"groupDescription\":\"example audit event\",\"vpcId\":\"vpc-a11a1111\"}",
	ResponseElements:  "{\"requestId\":\"add98acc-7c7d-44a0-829d-afef16381eff\",\"_return\":true,\"groupId\":\"sg-1234567890abcdefg\"}",
}

var eksEvent = EKSEvent{
	ClusterName: "example-eks",
	Kind:        "Event",
	APIVersion:  "audit.k8s.io/v1",
	Level:       "Metadata",
	AuditID:     "cec150e5-1af5-477a-a8b9-071846647f63",
	Stage:       "ResponseComplete",
	RequestURI:  "/api/v1/namespaces/default/configmaps/example?fieldManager=kubectl-edit\u0026fieldValidation=Strict",
	Verb:        "patch",
	User: EKSUser{
		Username: "example-user",
		Groups: []string{
			"5f8d568c-be13-48f3-8ac1-9349fa727bea",
			"d640f203-9a5c-4f56-9c7b-40da28bc2623",
		},
	},
	SourceIPs: []string{"1.2.3.4"},
	UserAgent: "kubectl/v1.30.0 (darwin/amd64) kubernetes/1234567",
	ObjectRef: EKSObjectRef{
		Resource:   "configmaps",
		Namespace:  "default",
		Name:       "example",
		APIVersion: "v1",
	},
	ResponseStatus: EKSResponseStatus{
		Metadata: map[string]string{},
		Code:     200,
	},
	RequestObject: EKSRequestObject{
		Data: map[string]string{
			"first":  "abc",
			"second": "def",
		},
	},
	ResponseObject: EKSResponseObject{
		Kind:       "ConfigMap",
		APIVersion: "v1",
		Data: map[string]string{
			"first":  "abc",
			"second": "def",
			"third":  "123",
		},
	},
	RequestReceivedTimestamp: time.Time{},
	StageTimestamp:           time.Time{},
	Annotations: map[string]string{
		"authorization.k8s.io/decision": "allow",
	},
}

func TestBuildCloudtrailMessage(t *testing.T) {
	t.Run("Build Cloudtrail Message", func(t *testing.T) {

		msg := buildCloudTrailMessage(cloudtrailEvent)

		assert.Equal(t, 4, len(msg.Blocks.BlockSet))
		assert.Equal(t, slack.MBTHeader, msg.Blocks.BlockSet[0].BlockType())
		assert.Equal(t, slack.MBTSection, msg.Blocks.BlockSet[1].BlockType())
		assert.Equal(t, slack.MBTSection, msg.Blocks.BlockSet[2].BlockType())
		assert.Equal(t, slack.MBTSection, msg.Blocks.BlockSet[3].BlockType())
	})
}

func TestBuildEKSMessage(t *testing.T) {
	t.Run("Build EKS Message", func(t *testing.T) {

		msg := buildEKSMessage(eksEvent)

		assert.Equal(t, 5, len(msg.Blocks.BlockSet))
		assert.Equal(t, slack.MBTHeader, msg.Blocks.BlockSet[0].BlockType())
		assert.Equal(t, slack.MBTSection, msg.Blocks.BlockSet[1].BlockType())
		assert.Equal(t, slack.MBTSection, msg.Blocks.BlockSet[2].BlockType())
		assert.Equal(t, slack.MBTSection, msg.Blocks.BlockSet[3].BlockType())
		assert.Equal(t, slack.MBTSection, msg.Blocks.BlockSet[4].BlockType())
	})
}
