package domain

// ReserveStatus represents current status of the accounts balance reservation
type ReserveStatus uint8

const (
	Reserved ReserveStatus = iota
	Accepted
	Declined
	UnknownReserve // добавлять новые статусы только в конец
)

func (status ReserveStatus) String() string {
	switch status {
	case Accepted:
		return "accepted"
	case Declined:
		return "declined"
	case Reserved:
		return "reserved"
	default:
		return "unknown"
	}
}

// TransactionStatus represents current transaction status of the reservation in accounting_report
type TransactionStatus uint8

const (
	Deposit TransactionStatus = iota
	Withdraw
	DeclinedTransaction
	UnknownTransaction
	AcceptedTransaction
)

func (status TransactionStatus) String() string {
	switch status {
	case Deposit:
		return "deposit"
	case Withdraw:
		return "withdraw"
	case DeclinedTransaction:
		return "declined"
	case AcceptedTransaction:
		return "accepted"
	default:
		return "unknown"
	}
}
