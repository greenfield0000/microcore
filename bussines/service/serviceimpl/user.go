package serviceimpl

import (
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type UserServiceImpl struct {
	repository *repository.Repository
	logger     *logrus.Logger
}

func NewUserService(repository *repository.Repository, logger *logrus.Logger) service.UserService {
	return &UserServiceImpl{logger: logger, repository: repository}
}

func (u *UserServiceImpl) GetUserInfo(id uint64) (*domains.User, error) {
	return u.repository.UserRepository.GetById(id)
}

func (u UserServiceImpl) Create(user domains.User) (domains.User, error) {
	return u.repository.UserRepository.Create(user)
}

func (u UserServiceImpl) Update(user *domains.User) error {
	return u.repository.UserRepository.Update(user)
}

func (u UserServiceImpl) All() ([]*domains.User, error) {
	return u.repository.UserRepository.All()
}

func (u *UserServiceImpl) IsExistById(id uint64) (bool, error) {
	return false, nil
}

func (u *UserServiceImpl) GetUserByAccountId(accountId uint64) (*domains.User, error) {
	return u.repository.UserRepository.GetUserByAccountId(accountId)
}

func (u *UserServiceImpl) DeleteById(id uint64) (bool, error) {
	return u.repository.UserRepository.DeleteById(id)
}
