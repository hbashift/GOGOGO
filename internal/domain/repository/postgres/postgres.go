package postgres

import (
	"HTTP-REST-API/internal/domain"
	"HTTP-REST-API/internal/domain/repository"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"time"
)

type postgresDb struct {
	postgres *sqlx.DB
}

const (
	HOST = "local"
	PORT = "8080"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// table creation

func InitPostgresDb(db *sqlx.DB) repository.Repository {

	path, _ := ioutil.ReadFile("configs/tables.sql")
	c := string(path)

	db.MustExec(c)

	return &postgresDb{postgres: db}
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetBalance returns error if you send wrong account_id
// else returns pointer to repository.Account object and nil error
func (db *postgresDb) GetBalance(accountId int) (*repository.Account, error) {
	account := repository.Account{}
	err := db.postgres.Get(&account, "SELECT * FROM account WHERE account_id=$1", accountId)

	if err != nil {
		account = repository.Account{}
		err = errors.New("there is no account with such account_id")
	}

	return &account, err
}

// AddToBalance returns status domain.Deposit if everything is okay and nil error
// else domain.DeclinedTransaction and error
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
			return domain.DeclinedTransaction, errors.New("database error")
		}

		db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2",
			uint32(amount)+account.Balance, accountId)

		return domain.Deposit, nil
	} else {
		db.postgres.MustExec("INSERT INTO account VALUES ($1, $2)", accountId, amount)

		return domain.Deposit, nil
	}
}

// checkBalance checks if amount is lesser or equals to the domain.Account.Balance
// returns domain.Reserved if okay else domain.Declined
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

// ReserveAmount scans through DataBase and returns domain.ReserveStatus if:
// 1) account exists
// 2) amount is lesser or equals to the account.Balance
// 3) if sum of all reservations referenced to the account.ID + amount if lesser than account.Balance
func (db *postgresDb) ReserveAmount(accountId, serviceId, orderId, amount int) (domain.ReserveStatus, error) {
	rows, _ := db.postgres.Query("SELECT EXISTS(SELECT * FROM reservation WHERE order_id=$1 AND (reservation_status=$2 OR reservation_status=$3 OR reservation_status=$4))",
		orderId, domain.Reserved.String(), domain.Declined.String(), domain.Accepted.String())

	var isExists bool

	for rows.Next() {
		if err := rows.Scan(&isExists); err != nil {
			log.Fatalln(err)
		}
	}

	if isExists {
		return domain.UnknownReserve, errors.New("reservation with such order_id is already exists")
	}

	rows, _ = db.postgres.Query("SELECT EXISTS(SELECT * FROM reservation WHERE order_id=$1 AND reservation_status=$2)",
		orderId, domain.Reserved.String())

	var reservationExists bool

	for rows.Next() {
		if err := rows.Scan(&reservationExists); err != nil {
			log.Fatalln(err)
		}
	}

	var status = domain.UnknownReserve

	if !reservationExists {
		// Checking if account with accountId exists
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
			return status, errors.New("")
		}

	} else {
		// Calculation all reservations with status "reserved"
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

// withdraw - auxiliary function, withdraws money from the account's balance
// returns domain.Withdraw
func (db *postgresDb) withdraw(accountId, amount int) (domain.TransactionStatus, error) {
	var status domain.TransactionStatus
	account := repository.Account{}
	err := db.postgres.Get(&account, "SELECT * FROM account WHERE account_id=$1", accountId)

	if err != nil {
		err = errors.New("there is no account with such account_id")
		return domain.DeclinedTransaction, err
	}

	db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2",
		account.Balance-uint32(amount), accountId)

	status = domain.Withdraw

	return status, err
}

// Admit returns domain.Withdraw if everything is okay
// if withdraw returns not nil err, Admit returns err and domain.DeclinedTransaction
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

		return status, err
	}

	db.postgres.MustExec("INSERT INTO accounting_report(service_id, order_id, amount, account_id, report_date, status) VALUES ($1, $2, $3, $4, $5, $6)",
		serviceId, orderId, amount, accountId, time.Now().Format("2006-01-02"), domain.AcceptedTransaction.String())

	return status, err
}

// deposit - auxiliary function, increases account.Balance if purchase was declined
func (db *postgresDb) deposit(accountId, amount int) (domain.TransactionStatus, error) {
	var status domain.TransactionStatus
	account := repository.Account{}
	err := db.postgres.Get(&account, "SELECT * FROM account WHERE account_id=$1", accountId)

	if err != nil {
		err = errors.New("there is no account with such account_id")
		return domain.DeclinedTransaction, err
	}

	db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2",
		account.Balance+uint32(amount), accountId)

	status = domain.Deposit

	return status, err
}

// DeclinePurchase returns domain.Deposit and nil error if purchase was declined
// and everything is okay, else domain.DeclinedTransaction and error
func (db *postgresDb) DeclinePurchase(accountId, orderId, serviceId int) (domain.TransactionStatus, error) {
	reserve := repository.Reservation{}
	err := db.postgres.Get(&reserve, "SELECT * FROM reservation WHERE reservation_status=$1 AND order_id=$2",
		domain.Accepted.String(), orderId)

	if err != nil {
		err = errors.New("there is no reservation with such order_id and status \"accepted\"")
		return domain.UnknownTransaction, err
	}

	if reserve.AccountId != accountId || reserve.ServiceId != serviceId {
		return domain.DeclinedTransaction, errors.New("bad params")
	}

	db.postgres.MustExec("UPDATE reservation SET reservation_status=$1 WHERE order_id=$2",
		domain.Declined.String(), orderId)

	status, err := db.deposit(accountId, int(reserve.Amount))

	if err != nil {
		db.postgres.MustExec("UPDATE reservation SET reservation_status=$1 WHERE order_id=$2",
			domain.Reserved.String(), orderId)

		return status, err
	}

	db.postgres.MustExec("UPDATE accounting_report SET status=$1 WHERE order_id=$2",
		domain.DeclinedTransaction.String(), orderId)

	return status, err
}

// TransferFromAccountToAccount returns domain.Deposit and nil error if accounts are existing and sender's balance
// is bigger than the amount of transfer, else domain.DeclinedTransaction and error
func (db *postgresDb) TransferFromAccountToAccount(accountId1, accountId2, amount int) (domain.TransactionStatus, error) {
	account1 := &repository.Account{}
	account2 := &repository.Account{}
	// Trying to find sender
	err := db.postgres.Get(account1, "SELECT * FROM account WHERE account_id=$1", accountId1)

	if err != nil {
		err = errors.New("there is no sender account with such account_id")
		return domain.DeclinedTransaction, err
	}
	// Trying to find receiver
	err = db.postgres.Get(account2, "SELECT * FROM account WHERE account_id=$1", accountId2)

	if err != nil {
		err = errors.New("there is no receiver account with such account_id")
		return domain.DeclinedTransaction, err
	}
	// If everything is okay, updating their balances
	db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2",
		account1.Balance-uint32(amount), accountId1)

	db.postgres.MustExec("UPDATE account SET balance=$1 WHERE account_id=$2",
		account2.Balance+uint32(amount), accountId2)

	return domain.Deposit, nil
}
