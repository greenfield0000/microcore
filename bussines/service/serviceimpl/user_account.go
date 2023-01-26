package serviceimpl

import (
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type UserAccountServiceImpl struct {
	repository *repository.Repository
	logger     *logrus.Logger
}

func NewUserAccountService(repository *repository.Repository, logger *logrus.Logger) service.UserAccountService {
	return &UserAccountServiceImpl{repository: repository, logger: logger}
}

func (ua UserAccountServiceImpl) Create(userAccount domains.UserAccount) (uint64, error) {
	return ua.repository.UserAccountRepository.Create(userAccount)
}

func (ua UserAccountServiceImpl) DeleteById(id uint64) (bool, error) {
	_, err := ua.repository.UserAccountRepository.DeleteById(id)
	return err == nil, err
}

func (ua UserAccountServiceImpl) GetUserAccountByAccountId(accountId uint64) (*domains.UserAccount, error) {
	return ua.repository.UserAccountRepository.GetByAccountId(accountId)
}
