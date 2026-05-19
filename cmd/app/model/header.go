package model

type Header struct {
	ClientID   string `reqHeader:"X-CLIENT-ID" json:"X-CLIENT-ID" validate:"required,max=64"`
	ExternalID string `reqHeader:"X-EXTERNAL-ID" json:"X-EXTERNAL-ID" validate:"required,max=64,alphanum"`
}
