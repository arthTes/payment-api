package account

type Request struct {
	DocumentNumber string `json:"document_number" binding:"required"`
}
