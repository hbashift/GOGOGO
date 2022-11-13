package repository

import (
	"HTTP-REST-API/internal/domain"
	_ "github.com/go-sql-driver/mysql"
)

type Repository interface {
	GetBalance(accountId int) (*Account, error)
	AddToBalance(accountId, amount int) (domain.TransactionStatus, error)
	ReserveAmount(accountId, serviceId, orderId, amount int) (domain.ReserveStatus, error)
	Admit(accountId, orderId, serviceId, amount int) (domain.TransactionStatus, error)
}
