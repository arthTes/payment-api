package domain

type Account struct {
	Id             string
	DocumentNumber string
}

func NewAccount(id, documentNumber string) Account {
	return Account{
		Id:             id,
		DocumentNumber: documentNumber,
	}
}
