package entities

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

func ReserveStatusFromString(val string) ReserveStatus {
	switch val {
	case "reserved":
		return Reserved
	case "accepted":
		return Accepted
	case "declined":
		return Declined
	default:
		return UnknownReserve
	}
}

type TransactionStatus uint8

const (
	Deposit TransactionStatus = iota
	Withdraw
	UnknownTransaction
)

func (status TransactionStatus) String() string {
	switch status {
	case Deposit:
		return "deposit"
	case Withdraw:
		return "withdraw"
	default:
		return "unknown"
	}
}

func TransactionStatusFromString(val string) TransactionStatus {
	switch val {
	case "deposit":
		return Deposit
	case "withdraw":
		return Withdraw
	default:
		return UnknownTransaction
	}
}
