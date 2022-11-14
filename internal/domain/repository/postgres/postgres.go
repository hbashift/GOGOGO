package postgres

import (
	"HTTP-REST-API/internal/domain"
	"HTTP-REST-API/internal/domain/repository"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
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
    reservation_id serial PRIMARY KEY,
    account_id int REFERENCES account(account_id),
    service_id int NOT NULL,
    order_id int NOT NULL UNIQUE,
    amount decimal NOT NULL,
    reserve_date date NOT NULL,
    reservation_status varchar(20) NOT NULL
);

CREATE TABLE accounting_report (
    report_id serial PRIMARY KEY,
    service_id int NOT NULL,
    order_id int NOT NULL UNIQUE,
    amount decimal NOT NULL,
    account_id int,
    report_date date                        
)`

// table creation

func InitPostgresDb(db *sqlx.DB) repository.Repository {
	// TODO Setup() создает БД и таблички либо sql файл
	// db.MustExec(schema)
	return &postgresDb{postgres: db}
}

func (db *postgresDb) GetBalance(accountId int) (*repository.Account, error) {
	account := repository.Account{}
	err := db.postgres.Get(&account, "SELECT * FROM account WHERE account_id=$1", accountId)

	if err != nil {
		account = repository.Account{}
		err = errors.New("there is no user with such account_id")
	}

	return &account, err
}

func (db *postgresDb) AddToBalance(accountId, amount int) (domain.TransactionStatus, error) {
	rows, _ := db.postgres.Query("SELECT EXISTS(SELECT * FROM account WHERE account_id=$1)", accountId)

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
			return domain.UnknownTransaction, errors.New("database error")
		}

		db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2",
			uint32(amount)+account.Balance, accountId)

		return domain.Deposit, nil
	} else {
		db.postgres.MustExec("INSERT INTO account VALUES ($1, $2)", accountId, amount)

		return domain.Deposit, nil
	}
}

func (db *postgresDb) checkBalance(accountId, amount int) domain.ReserveStatus {
	var status domain.ReserveStatus
	account, err := db.GetBalance(accountId)

	if err != nil {
		panic(err)
	}

	if uint32(amount) > account.Balance {
		status = domain.Declined
	} else {
		status = domain.Reserved
	}

	return status
}

func (db *postgresDb) ReserveAmount(accountId, serviceId, orderId, amount int) (domain.ReserveStatus, error) {
	rows, _ := db.postgres.Query("SELECT EXISTS(SELECT * FROM reservation WHERE order_id=$1 AND reservation_status=$2)",
		orderId, domain.Accepted.String())

	var isExists bool

	for rows.Next() {
		if err := rows.Scan(&isExists); err != nil {
			log.Fatalln(err)
		}
	}

	if isExists {
		return domain.UnknownReserve, errors.New("such order with order_id is already accepted")
	}

	rows, _ = db.postgres.Query("SELECT EXISTS(SELECT * FROM reservation WHERE account_id=$1 AND reservation_status=$2 AND order_id=$3)",
		accountId, domain.Reserved.String(), orderId)

	isExists = false

	for rows.Next() {
		if err := rows.Scan(&isExists); err != nil {
			log.Fatalln(err)
		}
	}

	var status = domain.UnknownReserve

	if !isExists {
		rows, _ := db.postgres.Queryx("SELECT EXISTS(SELECT * FROM account WHERE account_id=$1)", accountId)

		var accountExists bool

		for rows.Next() {
			if err := rows.Scan(&accountExists); err != nil {
				return domain.UnknownReserve, err
			}
		}

		if !accountExists {
			return domain.UnknownReserve, errors.New("there is no account with such account_id")
		}

		if status = db.checkBalance(accountId, amount); status == domain.Reserved {
			db.postgres.MustExec("INSERT INTO reservation(account_id, service_id, order_id, amount, reserve_date, reservation_status) VALUES ($1, $2, $3, $4, $5, $6)",
				accountId, serviceId, orderId, amount, time.Now().Format("2006-01-02"), status.String())
		} else {
			return status, errors.New("huyhuy")
		}
	} else {
		rows, _ := db.postgres.Queryx("SELECT EXISTS(SELECT * FROM reservation WHERE order_id=$1)", orderId)

		var orderExists bool

		for rows.Next() {
			if err := rows.Scan(&orderExists); err != nil {
				log.Fatalln(err)
			}
		}

		if orderExists {
			return domain.UnknownReserve, errors.New("reservation with such order_id already exists")
		}

		var reserves []repository.Reservation
		err := db.postgres.Select(&reserves, "SELECT * FROM reservation WHERE account_id=$1 AND reservation_status=$2",
			accountId, domain.Reserved.String())
		if err != nil {
			return domain.UnknownReserve, err
		}

		var expenses uint32

		for _, reserve := range reserves {
			expenses += reserve.Amount
		}

		account := repository.Account{}
		err = db.postgres.Get(&account, "SELECT * FROM account WHERE account_id=$1", accountId)
		if err != nil {
			panic(err)
		}

		if expenses+uint32(amount) > account.Balance {
			return domain.Declined, errors.New("not enough balance")
		} else {
			status = db.checkBalance(accountId, amount)

			db.postgres.MustExec("INSERT INTO reservation(account_id, service_id, order_id, amount, "+
				"reserve_date, reservation_status) VALUES ($1, $2, $3, $4, $5, $6)",
				accountId, serviceId, orderId, amount, time.Now().Format("2006-01-02"), status.String())
		}
	}
	return status, nil
}

func (db *postgresDb) withdraw(accountId, amount int) (domain.TransactionStatus, error) {
	var status domain.TransactionStatus
	account := repository.Account{}
	err := db.postgres.Get(&account, "SELECT * FROM account WHERE account_id=$1", accountId)

	if err != nil {
		err = errors.New("there is no user with such account_id")
	}

	db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2",
		account.Balance-uint32(amount), accountId)

	status = domain.Withdraw

	return status, err
}

func (db *postgresDb) Admit(accountId, orderId, serviceId, amount int) (domain.TransactionStatus, error) {
	reserve := repository.Reservation{}
	err := db.postgres.Get(&reserve, "SELECT * FROM reservation WHERE reservation_status=$1 AND order_id=$2",
		domain.Reserved.String(), orderId)

	if err != nil {
		err = errors.New("there is no reservation with such order_id and status \"reserved\"")
		return domain.UnknownTransaction, err
	}

	if reserve.Amount != uint32(amount) || reserve.AccountId != accountId || reserve.ServiceId != serviceId {
		return domain.UnknownTransaction, errors.New("bad params")
	}

	db.postgres.MustExec("UPDATE reservation SET reservation_status=$1 WHERE order_id=$2",
		domain.Accepted.String(), orderId)

	status, err := db.withdraw(accountId, amount)

	if err != nil {
		db.postgres.MustExec("UPDATE reservation SET reservation_status=$1 WHERE order_id=$2",
			domain.Reserved.String(), orderId)
		account := repository.Account{}

		err := db.postgres.Get(&account, "SELECT * FROM account WHERE account_id=$1", accountId)

		if err != nil {
			return domain.UnknownTransaction, err
		}

		db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2",
			account.Balance+uint32(amount), accountId)

		status = domain.UnknownTransaction

		return status, err
	}

	db.postgres.MustExec("INSERT INTO accounting_report(service_id, order_id, amount, account_id, report_date) VALUES ($1, $2, $3, $4, $5)",
		serviceId, orderId, amount, accountId, time.Now().Format("2006-01-02"))

	return status, err
}

func (db *postgresDb) DeclinePurchase(accountId, orderId, serviceId, amount int) (domain.TransactionStatus, error) {
	reserve := repository.Reservation{}
	err := db.postgres.Get(&reserve, "SELECT * FROM reservation WHERE reservation_status=$1 AND order_id=$2",
		domain.Reserved.String(), orderId)

	if err != nil {
		err = errors.New("there is no reservation with such order_id and status \"reserved\"")
		return domain.UnknownTransaction, err
	}

	if reserve.Amount != uint32(amount) || reserve.AccountId != accountId || reserve.ServiceId != serviceId {
		return domain.UnknownTransaction, errors.New("bad params")
	}

	db.postgres.MustExec("UPDATE reservation SET reservation_status=$1 WHERE order_id=$2",
		domain.Declined.String(), orderId)

	status, err := db.withdraw(accountId, amount)

	if err != nil {
		db.postgres.MustExec("UPDATE reservation SET reservation_status=$1 WHERE order_id=$2",
			domain.Reserved.String(), orderId)
		account := repository.Account{}

		err := db.postgres.Get(&account, "SELECT * FROM account WHERE account_id=$1", accountId)

		if err != nil {
			return domain.UnknownTransaction, err
		}

		db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2",
			account.Balance+uint32(amount), accountId)

		status = domain.UnknownTransaction

		return status, err
	}

	db.postgres.MustExec("INSERT INTO accounting_report(service_id, order_id, amount, account_id, report_date) VALUES ($1, $2, $3, $4, $5)",
		serviceId, orderId, amount, accountId, time.Now().Format("2006-01-02"))

	return status, err
}

func (db *postgresDb) TransferFromAccountToAccount(accountId1, accountId2 int) (account1, account2 domain.TransactionStatus,
	err error) {
	return 0, 0, err
}
