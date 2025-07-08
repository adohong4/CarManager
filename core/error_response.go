package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adohong4/carZone/utils"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// return error message
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Status, e.Message)
}

// create new ErrorResponse
func NewErrorResponse(message string, status int) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Status:  status,
	}
}

type ErrorResponseInterface interface {
	Error() string
	GetStatus() int
	GetMessage() string
}

// GetStatus trả về mã trạng thái
func (e *ErrorResponse) GetStatus() int {
	return e.Status
}

// GetMessage trả về thông điệp lỗi
func (e *ErrorResponse) GetMessage() string {
	return e.Message
}

type ConflictRequestError struct {
	*ErrorResponse
}

type BadRequestError struct {
	*ErrorResponse
}

type AuthFailureError struct {
	*ErrorResponse
}

type NotFoundError struct {
	*ErrorResponse
}

type ForbiddenError struct {
	*ErrorResponse
}

func NewConflictRequestError(message ...string) *ConflictRequestError {
	msg := utils.HTTPStatusMap[utils.Conflict].Reason
	if len(message) > 0 {
		msg = message[0]
	}
	return &ConflictRequestError{
		ErrorResponse: NewErrorResponse(msg, utils.Conflict),
	}
}

func NewBadRequestError(message ...string) *BadRequestError {
	msg := utils.HTTPStatusMap[utils.BadRequest].Reason
	if len(message) > 0 {
		msg = message[0]
	}
	return &BadRequestError{
		ErrorResponse: NewErrorResponse(msg, utils.BadRequest),
	}
}

func NewAuthFailureError(message ...string) *AuthFailureError {
	msg := utils.HTTPStatusMap[utils.Unauthorized].Reason
	if len(message) > 0 {
		msg = message[0]
	}
	return &AuthFailureError{
		ErrorResponse: NewErrorResponse(msg, utils.Unauthorized),
	}
}

func NewNotFoundError(message ...string) *NotFoundError {
	msg := utils.HTTPStatusMap[utils.NotFound].Reason
	if len(message) > 0 {
		msg = message[0]
	}
	return &NotFoundError{
		ErrorResponse: NewErrorResponse(msg, utils.NotFound),
	}
}

func NewForbiddenError(message ...string) *ForbiddenError {
	msg := utils.HTTPStatusMap[utils.Forbidden].Reason
	if len(message) > 0 {
		msg = message[0]
	}
	return &ForbiddenError{
		ErrorResponse: NewErrorResponse(msg, utils.Forbidden),
	}
}

// Send Error Response
func SendErrorResponse(w http.ResponseWriter, err *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.GetStatus())
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  err.GetStatus(),
		"message": err.GetMessage(),
	})
}
