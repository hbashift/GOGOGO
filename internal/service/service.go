package service

import (
	"HTTP-REST-API/internal/domain"
	"HTTP-REST-API/internal/domain/repository"
)

type Service struct {
	repository repository.Repository
}

func InitService(repository repository.Repository) *Service {
	if repository == nil {
		return nil
	}

	return &Service{repository: repository}
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

func (s *Service) Admit(accountId, orderId, serviceid, amount int) (domain.TransactionStatus, error) {
	status, err := s.repository.Admit(accountId, orderId, serviceid, amount)

	if err != nil {
		return status, err
	}

	return status, nil
}
