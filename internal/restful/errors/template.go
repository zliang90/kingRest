package errors

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

type (
	Params map[string]interface{}

	errorTemplate struct {
		Code             int64  `yaml:"code"`
		Message          string `yaml:"message"`
		DeveloperMessage string `yaml:"developer_message"`
	}
)

var templates map[string]errorTemplate

func LoadMessages(file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	templates = map[string]errorTemplate{}
	return yaml.Unmarshal(bytes, &templates)
}

func NewAPIError(code string, params Params) *APIError {
	err := &APIError{
		Message: code,
	}

	if template, ok := templates[code]; ok {
		err.Code = template.getErrorCode()
		err.Message = template.getMessage(params)

		if Env != "prod" {
			err.DeveloperMessage = template.getDeveloperMessage(params)
		}
	}

	return err
}

func (e errorTemplate) getMessage(params Params) string {
	return replacePlaceholders(e.Message, params)
}

func (e errorTemplate) getDeveloperMessage(params Params) string {
	return replacePlaceholders(e.DeveloperMessage, params)
}

func (e errorTemplate) getErrorCode() int64 {
	return e.Code
}

func replacePlaceholders(message string, params Params) string {
	if len(message) == 0 {
		return ""
	}
	if params != nil {
		for key, value := range params {
			message = strings.Replace(message, "{"+key+"}", fmt.Sprint(value), -1)
		}
	}
	return message
}
