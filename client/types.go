package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aas-core-works/aas-core3.1-golang/types"
)

type BasyxPagedResultMetadata struct {
	Cursor string `json:"cursor"`
}

type BasyxPagedResultRaw struct {
	Metadata BasyxPagedResultMetadata `json:"paging_metadata"`
	Result   []json.RawMessage        `json:"result"`
}

type BasyxPagedResult[T types.IClass] struct {
	Metadata BasyxPagedResultMetadata
	Result   []T
}

type BasyxErrorResult struct {
	Messages []BasyxErrorResultMessage `json:"messages"`
}

type BasyxErrorResultMessage struct {
	Timestamp     string `json:"timestamp"`
	Text          string `json:"text"`
	Code          string `json:"code"`
	CorrelationID string `json:"correlationId"`
	MessageType   string `json:"messageType"`
}

func (r BasyxErrorResult) Error() string {
	if len(r.Messages) == 0 {
		return "BaSyx returned an unknown error"
	}

	var parts []string
	for _, msg := range r.Messages {
		var fields []string

		if msg.Code != "" {
			fields = append(fields, fmt.Sprintf("code=%s", msg.Code))
		}
		if msg.Text != "" {
			fields = append(fields, msg.Text)
		}
		if msg.MessageType != "" {
			fields = append(fields, fmt.Sprintf("type=%s", msg.MessageType))
		}
		if msg.CorrelationID != "" {
			fields = append(fields, fmt.Sprintf("correlationId=%s", msg.CorrelationID))
		}
		if msg.Timestamp != "" {
			fields = append(fields, fmt.Sprintf("timestamp=%s", msg.Timestamp))
		}

		parts = append(parts, strings.Join(fields, ", "))
	}

	return strings.Join(parts, "; ")
}
