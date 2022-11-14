package domain

type Account struct {
	Id      int    `json:"account_id"`
	Balance uint32 `json:"balance"`
}

type AccountDto struct {
	Id           int    `json:"account_id"`
	BalanceAdded uint32 `json:"balance_added"`
}

type ReservationDto struct {
	AccountId int    `json:"account_id"`
	ServiceId int    `json:"service_id"`
	OrderId   int    `json:"order_id"`
	Amount    uint32 `json:"amount"`
}

type DeclineDto struct {
	AccountId int `json:"account_id"`
	OrderId   int `json:"order_id"`
	ServiceId int `json:"service_id"`
}

type Transfer struct {
	Sender   int    `json:"sender_id"`
	Receiver int    `json:"receiver_id"`
	Amount   uint32 `json:"amount"`
}
