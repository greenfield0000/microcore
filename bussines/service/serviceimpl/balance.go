package serviceimpl

import (
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type BalanceServiceImpl struct {
	repository *repository.Repository
}

func NewBalanceService(repository *repository.Repository, logger *logrus.Logger) service.BalanceService {
	return &BalanceServiceImpl{
		repository: repository,
	}
}

func (b BalanceServiceImpl) Create(balance *domains.Balance) error {
	return b.repository.BalanceRepository.Create(balance)
}

func (b BalanceServiceImpl) GetBalance(accountId uint64) (*domains.Balance, error) {
	return b.repository.GetBalance(accountId)
}
