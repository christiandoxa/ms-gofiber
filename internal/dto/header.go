package dto

type RequestHeader struct {
	XPartnerID  string `reqHeader:"X-PARTNER-ID" json:"X-PARTNER-ID" validate:"required,alphanum,max=36"`
	ChannelID   string `reqHeader:"CHANNEL-ID" json:"CHANNEL-ID" validate:"required,alphanum,max=5"`
	XExternalID string `reqHeader:"X-EXTERNAL-ID" json:"X-EXTERNAL-ID" validate:"required,numeric,max=36"`
}
