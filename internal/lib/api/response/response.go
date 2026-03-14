package response

import (
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

func OK() Response {
	return Response{Status: StatusOK}
}

func Error(err string) Response {
	return Response{Status: StatusError, Error: err}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, err.Field()+" is required")
		case "url":
			errMsg = append(errMsg, err.Field()+" must be a valid URL")
		default:
			errMsg = append(errMsg, err.Field()+" is invalid")
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, "; "),
	}
}
