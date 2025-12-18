package dto

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type CreateAccountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
