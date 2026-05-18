package model

type Header struct {
	ClientID string `reqHeader:"X-CLIENT-ID" json:"X-CLIENT-ID" validate:"required,max=64"`
}
