package domain

import (
	"strconv"
	"time"
)

type Account struct {
	Id      int    `json:"account_id"`
	Balance uint32 `json:"balance"`
}

type AccountDto struct {
	Id           int    `json:"account_id"`
	BalanceAdded uint32 `json:"balance_added"`
}

func (a Account) String() string {
	return "Id:" + strconv.FormatInt(int64(a.Id), 10) +
		"\n Balance" + strconv.FormatUint(uint64(a.Balance), 10)
}

type Service struct {
	Id     int    `json:"service_id"`
	Amount uint32 `json:"amount"`
	Name   string `json:"service_name"`
}

type Reservation struct {
	Id        int    `json:"reservation_id"`
	AccountId int    `json:"account_id"`
	ServiceId int    `json:"service_id"`
	OrderId   int    `json:"order_id"`
	Amount    uint32 `json:"amount"`
	Date      time.Time
	Status    ReserveStatus
}

type ReservationDto struct {
	AccountId int    `json:"account_id"`
	ServiceId int    `json:"service_id"`
	OrderId   int    `json:"order_id"`
	Amount    uint32 `json:"amount"`
}

type ReportDto struct {
	AccountId int    `json:"account_id"`
	OrderId   int    `json:"order_id"`
	ServiceId int    `json:"service_id"`
	Amount    uint32 `json:"amount"`
}
