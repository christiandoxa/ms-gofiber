package model

type ResponseCode struct {
	Endpoint              string
	OriginResponseCode    string
	OriginResponseMessage string
	ResponseCode          string
	ResponseMessage       string
	StatusCode            int
}
