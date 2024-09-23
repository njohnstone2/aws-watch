package main

import "time"

type Event struct {
	EventSource string
	EKS         *EKSEvent
	Cloudtrail  *CloudtrailEvent
}

type EKSEvent struct {
	ClusterName              string
	Kind                     string            `json:"kind"`
	APIVersion               string            `json:"apiVersion"`
	Level                    string            `json:"level"`
	AuditID                  string            `json:"auditID"`
	Stage                    string            `json:"stage"`
	RequestURI               string            `json:"requestURI"`
	Verb                     string            `json:"verb"`
	User                     EKSUser           `json:"user"`
	SourceIPs                []string          `json:"sourceIPs"`
	UserAgent                string            `json:"userAgent"`
	ObjectRef                EKSObjectRef      `json:"objectRef"`
	ResponseStatus           EKSResponseStatus `json:"responseStatus"`
	RequestObject            EKSRequestObject  `json:"requestObject"`
	ResponseObject           EKSResponseObject `json:"responseObject"`
	RequestReceivedTimestamp time.Time         `json:"requestReceivedTimestamp"`
	StageTimestamp           time.Time         `json:"stageTimestamp"`
	Annotations              map[string]string `json:"annotations"`
}

type EKSUser struct {
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
}

type EKSObjectRef struct {
	Resource   string `json:"resource"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	APIVersion string `json:"apiVersion"`
}

type EKSResponseStatus struct {
	Metadata map[string]string `json:"metadata"`
	Code     int               `json:"code"`
}

type EKSRequestObject struct {
	Data map[string]string `json:"data"`
}

type EKSResponseObject struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		Name              string            `json:"name"`
		Namespace         string            `json:"namespace"`
		UID               string            `json:"uid"`
		ResourceVersion   string            `json:"resourceVersion"`
		CreationTimestamp time.Time         `json:"creationTimestamp"`
		Labels            map[string]string `json:"labels"`
		Annotations       map[string]string `json:"annotations"`
	} `json:"metadata"`
	Data map[string]string `json:"data"`
}
