package serviceimpl

import (
	"errors"
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type MarketServiceImpl struct {
	repository *repository.Repository
	logger     *logrus.Logger
}

func NewMarketService(repository *repository.Repository, logger *logrus.Logger) service.MarketService {
	return &MarketServiceImpl{logger: logger, repository: repository}
}

func (m MarketServiceImpl) GetMarketInfo(id uint64) (*domains.Market, error) {
	return m.repository.MarketRepository.GetById(id)
}

func (m MarketServiceImpl) Create(market domains.Market) (uint64, error) {
	exist, err := m.IsExistBySysName(*market.Sysname)
	if err != nil {
		m.logger.Error(err)
		return 0, err
	}
	if exist {
		return 0, errors.New("Market already exist")
	}
	return m.repository.MarketRepository.Create(market)
}

func (m MarketServiceImpl) All() ([]*domains.Market, error) {
	return m.repository.MarketRepository.All()
}

func (m MarketServiceImpl) IsExistBySysName(sysName string) (bool, error) {
	return m.repository.MarketRepository.IsExistBySysName(sysName)
}

func (m MarketServiceImpl) IsExistById(id uint64) (bool, error) {
	return m.repository.MarketRepository.IsExistById(id)
}

func (m MarketServiceImpl) GetMarketsByAccountId(accountId uint64) (map[string][]*domains.Market, error) {
	markets, err := m.repository.MarketRepository.GetMarketsByAccountId(accountId)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}

	mapMarket := make(map[string][]*domains.Market, 3)
	mapMarket["IN_PROCESS"] = make([]*domains.Market, 0)
	mapMarket["ACQUIRED"] = make([]*domains.Market, 0)
	mapMarket["NEW"] = make([]*domains.Market, 0)

	for i := range markets {
		m := markets[i]
		if m.State == nil || m.State.Sysname == nil {
			if acMarket := mapMarket["NEW"]; !contains(acMarket, m) {
				mapMarket["NEW"] = append(mapMarket["NEW"], m)
			}
		} else {
			switch *m.State.Sysname {
			case "ACQUIRED":
				if acMarket := mapMarket["ACQUIRED"]; !contains(acMarket, m) {
					mapMarket["ACQUIRED"] = append(mapMarket["ACQUIRED"], m)
				}

			case "IN_PROCESS":
				if acMarket := mapMarket["IN_PROCESS"]; !contains(acMarket, m) {
					mapMarket["IN_PROCESS"] = append(mapMarket["IN_PROCESS"], m)
				}
			case "NEW":
			case "CANCELED":
				if acMarket := mapMarket["NEW"]; !contains(acMarket, m) {
					mapMarket["NEW"] = append(mapMarket["NEW"], m)
				}
			}
		}
	}

	return mapMarket, err
}

func contains(list []*domains.Market, item *domains.Market) bool {
	if list == nil || item == nil {
		return false
	}

	for _, v := range list {
		if v == item || *v.Id == *item.Id {
			return true
		}
	}
	return false
}
