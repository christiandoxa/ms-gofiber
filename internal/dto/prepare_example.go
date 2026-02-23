package dto

type PrepareExampleRequest struct {
	TerminalType      string   `json:"terminalType" validate:"required,terminal_type"`
	OsType            string   `json:"osType" validate:"omitempty,oneof=IOS ANDROID OTHER"`
	OsVersion         string   `json:"osVersion" validate:"omitempty,max=32"`
	GrantType         string   `json:"grantType" validate:"required,grant_type"`
	PaymentMethodType string   `json:"paymentMethodType" validate:"required,payment_method_type"`
	Scope             []string `json:"scope" validate:"required,authorization_scope"`
	TransactionTime   string   `json:"transactionTime" validate:"required,time_rule"`
	MerchantName      string   `json:"merchantName" validate:"required,min=3,max=100,alphanum_with_space"`
}
