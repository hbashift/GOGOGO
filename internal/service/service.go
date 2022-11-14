package service

import (
	"HTTP-REST-API/internal/domain"
	"HTTP-REST-API/internal/domain/repository"
	"errors"
)

type Service struct {
	repository repository.Repository
}

func InitService(repository repository.Repository) (*Service, error) {
	if repository == nil {
		return nil, errors.New("empty repository")
	}

	return &Service{
		repository: repository,
	}, nil
}

func (s *Service) GetBalance(id int) (*domain.Account, error) {
	account, err := s.repository.GetBalance(id)

	if err != nil {
		return nil, err
	}

	result := domain.Account{
		Id:      account.Id,
		Balance: account.Balance,
	}

	return &result, nil
}

func (s *Service) AddToBalance(accountId, amount int) error {
	_, err := s.repository.AddToBalance(accountId, amount)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ReserveAmount(accountId, serviceId, orderId, amount int) (domain.ReserveStatus, error) {
	status, err := s.repository.ReserveAmount(accountId, serviceId, orderId, amount)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (s *Service) Admit(accountId, orderId, serviceId, amount int) (domain.TransactionStatus, error) {
	status, err := s.repository.Admit(accountId, orderId, serviceId, amount)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (s *Service) Decline(accountId, orderId, serviceId int) (domain.TransactionStatus, error) {
	status, err := s.repository.DeclinePurchase(accountId, orderId, serviceId)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (s *Service) Transfer(accountId1, accountId2, amount int) (domain.TransactionStatus, error) {
	status, err := s.repository.TransferFromAccountToAccount(accountId1, accountId2, amount)

	if err != nil {
		return status, err
	}

	return status, nil
}
