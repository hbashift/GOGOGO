package repository

import (
	"time"
)

type Account struct {
	Id      int    `db:"account_id"`
	Balance uint32 `db:"balance"`
}

type Service struct {
	Id     int    `db:"service_id"`
	Amount uint32 `db:"amount"`
	Name   string `db:"service_name"`
}

type Reservation struct {
	Id        int       `db:"reservation_id"`
	AccountId int       `db:"account_id"`
	ServiceId int       `db:"service_id"`
	OrderId   int       `db:"order_id"`
	Amount    uint32    `db:"amount"`
	Date      time.Time `db:"reserve_date"`
	Status    string    `db:"reservation_status"`
}
