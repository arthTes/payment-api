package operation

type Type int

const (
	CASH_PURCHASES Type = iota
	INSTALLMENT_PURCHASES
	WITHDRAW
	PAYMENT
)

func (t Type) String() string {
	return [...]string{"CASH_PURCHASES", "INSTALLMENT_PURCHASES", "WITHDRAW", "PAYMENT"}[t]
}

func (t Type) Index() int {
	return int(t)
}

func (t Type) IsValid() bool {
	if t.Index() > 3 || t.Index() < 0 {
		return false
	}

	return true
}
