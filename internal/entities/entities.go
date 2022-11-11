package entities

import (
	"time"
)

type Account struct {
	Id      int    `json:"account_id"`
	Balance uint32 `json:"balance"`
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

/*var Accounts = []Account{
	{Id: 1, Balance: decimal.New(100, 2)},
	{Id: 2, Balance: decimal.New(150, 0)},
	{Id: 3, Balance: decimal.New(100, 0)},
	{Id: 4, Balance: decimal.New(1000, 2)},
}

var Services = []Service{
	{Id: 1, Amount: decimal.New(200, 2)},
	{Id: 2, Amount: decimal.New(250, 0)},
	{Id: 3, Amount: decimal.New(200, 0)},
	{Id: 4, Amount: decimal.New(2000, 2)},
}

var Reservations = []Reservation{
	{Id: 1, AccountId: 1, ServiceId: 1, OrderId: 10, Amount: decimal.New(300, 2), Date: time.Now().Local(), Status: Accepted},
	{Id: 2, AccountId: 2, ServiceId: 2, OrderId: 20, Amount: decimal.New(350, 0), Date: time.Now().Local(), Status: Accepted},
	{Id: 3, AccountId: 3, ServiceId: 3, OrderId: 30, Amount: decimal.New(300, 0), Date: time.Now().Local(), Status: Accepted},
	{Id: 4, AccountId: 4, ServiceId: 4, OrderId: 40, Amount: decimal.New(3000, 2), Date: time.Now().Local(), Status: Accepted},
}*/
