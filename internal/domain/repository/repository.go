package repository

import (
	"HTTP-REST-API/internal/domain"
	_ "github.com/go-sql-driver/mysql"
)

type Repository interface {
	GetBalance(accountId int) (domain.Account, error)
	AddToBalance(accountId, amount int) (domain.TransactionStatus, error)
	ReserveAmount(accountId, serviceId, orderId, amount int) (domain.ReserveStatus, error)
	Withdraw(accountId, amount int) (domain.TransactionStatus, error)
	Admit(accountId, orderId, serviceId, amount int) (domain.TransactionStatus, error)
}
