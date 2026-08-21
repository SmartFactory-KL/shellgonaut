package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aas-core-works/aas-core3.1-golang/types"
)

// ------------------ Paged Results ------------
type PagedResultMetadata struct {
	Cursor string `json:"cursor"`
}

type PagedResultRaw struct {
	Metadata PagedResultMetadata `json:"paging_metadata"`
	Result   []json.RawMessage   `json:"result"`
}

type PagedResult[T types.IClass] struct {
	Metadata PagedResultMetadata
	Result   []T
}

type PagedStringResult struct {
	Metadata PagedResultMetadata
	Result   []string
}

// -------------- Error Results ----------------------

type ErrorResult struct {
	Messages []ErrorResultMessage `json:"messages"`
}

type ErrorResultMessage struct {
	Timestamp     string `json:"timestamp"`
	Text          string `json:"text"`
	Code          string `json:"code"`
	CorrelationID string `json:"correlationId"`
	MessageType   string `json:"messageType"`
}

func (r ErrorResult) Error() string {
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

// --------------------- Operation ----------------

type OperationRequest struct {
	InputArguments    []types.IOperationVariable
	InoutputArguments []types.IOperationVariable
}

type OperationRequestJSON struct {
	InputArguments    []json.RawMessage `json:"inputArguments"`
	InoutputArguments []json.RawMessage `json:"inoutputArguments"`
}

type OperationResult struct {
	Success           bool
	OutputArguments   []types.IOperationVariable
	InoutputArguments []types.IOperationVariable
}

type OperationResultJSON struct {
	Success           bool              `json:"success"`
	OutputArguments   []json.RawMessage `json:"outputArguments"`
	InoutputArguments []json.RawMessage `json:"inoutputArguments"`
}
