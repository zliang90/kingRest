package errors

type APIError struct {
	RequestId        string      `json:"request_id,omitempty"`
	Code             int64       `json:"code"`
	Message          string      `json:"message,omitempty"`
	DeveloperMessage string      `json:"developer_message,omitempty"`
	Details          interface{} `json:"details,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}
