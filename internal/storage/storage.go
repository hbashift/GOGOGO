package storage

import (
	"HTTP-REST-API/internal/entities"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

func CreateDataBase() {
	// TODO creating DB
}

type Storage interface {
	GetBalance(userId int) (int, error)
	AddToBalance(userId, amount int) (entities.TransactionStatus, error)
	ReserveAmount(userId, serviceId, orderId, amount int) (entities.ReserveStatus, error)
	WithDraw(userId, amount int) (entities.TransactionStatus, error)
	Admit(userId, orderId, serviceId, amount int) (entities.TransactionStatus, error)
}

type postgresDb struct {
	accountDb     *AccountDb
	reserveDb     *ReserveDb
	transactionDb *TransactionDb
}

type AccountDb struct {
	db   *sql.DB
	mute sync.Mutex
}

type ReserveDb struct {
	db   *sql.DB
	mute sync.Mutex
}

type TransactionDb struct {
	db   *sql.DB
	mute sync.Mutex
}
