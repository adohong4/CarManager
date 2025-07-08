package core

import (
	"encoding/json"
	"net/http"

	"github.com/adohong4/carZone/utils"
)

type SuccessResponse struct {
	Status   int         `json:"status"`
	Message  string      `json:"message"`
	Metadata interface{} `json:"metadata,omitempty"`
}

func NewSuccessResponse(message string, metadata interface{}) *SuccessResponse {
	return &SuccessResponse{
		Status:   utils.OK,
		Message:  message,
		Metadata: metadata,
	}
}

func (s *SuccessResponse) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s.Status)
	json.NewEncoder(w).Encode(s)
}

type OK struct {
	*SuccessResponse
}

type CREATED struct {
	*SuccessResponse
}

func NewOK(message string, metadata interface{}) *OK {
	if message == "" {
		message = utils.HTTPStatusMap[utils.OK].Reason
	}
	return &OK{
		SuccessResponse: &SuccessResponse{
			Status:   utils.OK,
			Message:  message,
			Metadata: metadata,
		},
	}
}

func NewCREATED(message string, metadata interface{}) *CREATED {
	if message == "" {
		message = utils.HTTPStatusMap[utils.Created].Reason
	}
	return &CREATED{
		SuccessResponse: &SuccessResponse{
			Status:   utils.Created,
			Message:  message,
			Metadata: metadata,
		},
	}
}
