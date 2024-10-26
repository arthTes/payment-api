package operation

type Type int

const (
	CASH_PURCHASES Type = iota + 1
	INSTALLMENT_PURCHASES
	WITHDRAW
	PAYMENT
)

func (t Type) String() string {
	return [...]string{"CASH_PURCHASES", "INSTALLMENT_PURCHASES", "WITHDRAW", "PAYMENT"}[t-1]
}

func (t Type) Index() int {
	return int(t)
}

func (t Type) IsValid() bool {
	if t.Index() > 4 || t.Index() < 1 {
		return false
	}

	return true
}
