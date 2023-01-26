package serviceimpl

import (
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type BalanceRobotServiceImpl struct {
	repository *repository.Repository
	logger     *logrus.Logger
}

func NewBalanceRobotService(repository *repository.Repository, logger *logrus.Logger) service.BalanceRobotService {
	return &BalanceRobotServiceImpl{
		logger:     logger,
		repository: repository,
	}
}

func (b BalanceRobotServiceImpl) Create(balance *domains.BalanceRobot) error {
	return b.repository.BalanceRobotRepository.Create(balance)
}
