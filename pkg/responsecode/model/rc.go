package rcmodel

type ResponseCode struct {
	StatusCode      int    `json:"-"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	ResponseDesc    string `json:"responseDesc,omitempty"`
}

func NewResponseCode(statusCode int, message string, responseCode string) *ResponseCode {
	return &ResponseCode{
		StatusCode:      statusCode,
		ResponseCode:    responseCode,
		ResponseMessage: message,
	}
}

func (r *ResponseCode) Error() string {
	if r.ResponseMessage != "" {
		return r.ResponseMessage
	}
	return r.ResponseDesc
}
