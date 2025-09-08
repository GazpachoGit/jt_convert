package tc_client

import (
	"errors"
	"fmt"
)

type GenericResponseBody struct {
	ServiceData   interface{}        `json:"ServiceData"`
	PartialErrors []ParialErrorsData `json:"partialErrors"`
}

type ParialErrorsData struct {
	ErrorValues []ErrorValue `json:"errorValues"`
}

type ErrorValue struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Level   string `json:"level"`
}

func (e GenericResponseBody) GetPartialErrors() error {
	if len(e.PartialErrors) > 0 {
		errMsg := "TC Partial errors: "
		for _, partialErr := range e.PartialErrors {
			for _, errVal := range partialErr.ErrorValues {
				errMsg += fmt.Sprintf("[Code: %s, Level: %s, Message: %s] ", errVal.Code, errVal.Level, errVal.Message)
			}
		}
		return errors.New(errMsg)
	}
	return nil
}
