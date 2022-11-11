package storage

import "HTTP-REST-API/internal/entities"

// init func for account table in database

func AccountDbInit(addr string) (AccountDb, error) {
	return AccountDb{}, nil
}

// init func for reservation table in database

func ReserveDbInit(addr string) (ReserveDb, error) {
	return ReserveDb{}, nil
}

func (db *postgresDb) GetBalance(userId int) (int, error) {
	// TODO SQL request
	return 0, nil
}

func (db *postgresDb) AddToBalance(userId, amount int) (entities.TransactionStatus, error) {
	return entities.UnknownTransaction, nil
}

func (db *postgresDb) ReserveAmount(userId, serviceId, orderId, amount int) (entities.ReserveStatus, error) {
	return entities.UnknownReserve, nil
}

func (db *postgresDb) WithDraw(userId, amount int) (entities.TransactionStatus, error) {
	return entities.UnknownTransaction, nil
}

func (db *postgresDb) Admit(userId, orderId, serviceId, amount int) (entities.TransactionStatus, error) {
	return entities.UnknownTransaction, nil
}
