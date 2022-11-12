package postgres

import (
	"HTTP-REST-API/internal/domain"
	"HTTP-REST-API/internal/domain/repository"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type postgresDb struct {
	postgres *sqlx.DB
}

var schema = `
CREATE TABLE account (
    account_id int PRIMARY KEY,
    balance decimal
);

CREATE TABLE reservation (
    reservation_id int PRIMARY KEY,
    account_id int REFERENCES account(account_id),
    service_id int NOT NULL,
    order_id int NOT NULL UNIQUE,
    amount int NOT NULL,
    reserve_date date NOT NULL,
    reservation_status varchar(15) NOT NULL
)`

// table creation

func InitPostgresDb(db *sqlx.DB) repository.Repository {
	// TODO Setup() создает БД и таблички либо sql файл
	db.MustExec(schema)
	fmt.Println("Tables are created")
	return &postgresDb{postgres: db}
}

func (db *postgresDb) GetBalance(accountId int) (*repository.Account, error) {
	account := repository.Account{}
	err := db.postgres.Get(&account, "SELECT * FROM account WHERE account_id=$1", accountId)
	if err != nil {
		fmt.Errorf("no such account with such Id")
	}
	return &account, err
}

func (db *postgresDb) AddToBalance(accountId, amount int) (domain.TransactionStatus, error) {

	rows, _ := db.postgres.Query("SELECT EXISTS(SELECT 1 FROM account WHERE account_id=$1)", accountId)

	var isExists bool
	for rows.Next() {
		if err := rows.Scan(&isExists); err != nil {
			log.Fatalln(err)
		}
	}

	if isExists {
		account := repository.Account{}
		err := db.postgres.Get(&account, "SELECT * FROM account WHERE account_id=$1", accountId)
		if err != nil {
			return domain.UnknownTransaction, errors.New("hehe poopoo")
		}
		db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2",
			uint32(amount)+account.Balance, accountId)

		return domain.Deposit, nil
	} else {
		db.postgres.MustExec("INSERT INTO account VALUES ($1, $2)", accountId, amount)

		return domain.Deposit, nil
	}

	/*if _, err := db.postgres.Exec("SELECT EXISTS(SELECT * FROM account WHERE account_id=$1)", accountId); err != nil {
		balance := 100
		db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2", amount+balance, accountId)
		return domain.Deposit, err
	}*/
	// TODO обработка ошибок и правильный INSERT
	/*	err := db.postgres.MustExec("INSERT INTO account VALUES ($1, $2)", accountId, amount)
		if err != nil {
			errors.New("could insert new account")
		}
		return domain.Deposit, nil*/
}

func (db *postgresDb) ReserveAmount(accountId, serviceId, orderId, amount int) (domain.ReserveStatus, error) {
	return domain.UnknownReserve, nil
}

func (db *postgresDb) Withdraw(accountId, amount int) (domain.TransactionStatus, error) {
	return domain.UnknownTransaction, nil
}

func (db *postgresDb) Admit(userId, orderId, serviceId, amount int) (domain.TransactionStatus, error) {
	return domain.UnknownTransaction, nil
}
