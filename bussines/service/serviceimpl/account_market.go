package serviceimpl

import (
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type AccountMarketServiceImpl struct {
	repository *repository.Repository
	logger     *logrus.Logger
}

func NewAccountMarketService(repository *repository.Repository, logger *logrus.Logger) service.AccountMarketService {
	return &AccountMarketServiceImpl{logger: logger, repository: repository}
}

func (aa AccountMarketServiceImpl) Create(accountMarket domains.AccountMarket) (uint64, error) {
	return aa.repository.AccountMarketRepository.Create(accountMarket)
}

func (aa AccountMarketServiceImpl) DeleteById(id uint64) (bool, error) {
	return aa.repository.AccountMarketRepository.DeleteById(id)
}

func (aa AccountMarketServiceImpl) IsMarketInAccount(marketId uint64, accountId uint64) (bool, error) {
	return aa.repository.AccountMarketRepository.IsMarketInAccount(marketId, accountId)
}
