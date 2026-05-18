package response

type Envelope struct {
	Status  string            `json:"status"`
	Message string            `json:"message,omitempty"`
	Data    any               `json:"data,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func Success(data any) Envelope {
	return Envelope{
		Status: "success",
		Data:   data,
	}
}

func Error(message string, fields map[string]string) Envelope {
	return Envelope{
		Status:  "error",
		Message: message,
		Fields:  fields,
	}
}
