package serviceimpl

import (
	"errors"
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type AccountServiceImpl struct {
	logger *logrus.Logger

	repository *repository.Repository
}

func NewAccountService(repository *repository.Repository, logger *logrus.Logger) service.AccountService {
	return &AccountServiceImpl{
		logger:     logger,
		repository: repository,
	}
}

func (s *AccountServiceImpl) Create(account domains.Account) (uint64, error) {
	return s.repository.AccountRepository.Create(account)
}

func (s *AccountServiceImpl) GetAccountInfo(id uint64) (*domains.Account, error) {
	return s.repository.AccountRepository.GetById(id)
}

func (s *AccountServiceImpl) All() ([]*domains.Account, error) {
	return s.repository.AccountRepository.All()
}

func (s *AccountServiceImpl) IsExistById(id uint64) (bool, error) {
	return s.repository.AccountRepository.IsExist(id)
}

func (s *AccountServiceImpl) DeleteById(id uint64) (bool, error) {
	return s.repository.AccountRepository.DeleteById(id)
}

func (s *AccountServiceImpl) AddAchievement(accountId uint64, achievementId uint64) (id uint64, err error) {
	//isAccountExist, _ := s.Service.AccountService.IsExistById(accountId)
	//if !isAccountExist {
	//	return 0, errors.New("Не удалось добавить достижение")
	//}
	//
	//isAchievementInAccount, _ := s.Service.AccountAchievementService.IsAchievementInAccount(achievementId, accountId)
	//if isAchievementInAccount {
	//	return 0, errors.New("У данного аккаунта уже есть такое достижение")
	//}
	//
	//return s.repository.AccountAchievementRepository.Create(domains.AccountAchievement{
	//	AccountId:     &accountId,
	//	AchievementId: &achievementId,
	//})
	return 0, errors.New("Переделать")
}

func (s *AccountServiceImpl) AddMarket(accountId uint64, marketId uint64) (id uint64, err error) {
	//isAccountExist, _ := s.Service.AccountService.IsExistById(accountId)
	//if !isAccountExist {
	//	return 0, errors.New("Не удалось добавить позицию магазина")
	//}
	//
	//isMarketExist, _ := s.Service.AccountMarketService.IsMarketInAccount(marketId, accountId)
	//if isMarketExist {
	//	return 0, errors.New("У данного аккаунта уже есть такая позиция магазина")
	//}
	//
	//isMarketInAccount, _ := s.Service.AccountMarketService.IsMarketInAccount(marketId, accountId)
	//if isMarketInAccount {
	//	return 0, errors.New("Данная позиция уже существует")
	//}
	//return s.repository.AccountMarketRepository.Create(domains.AccountMarket{
	//	AccountId: &accountId,
	//	MarketId:  &marketId,
	//})
	return 0, errors.New("Переделать")
}

func (s *AccountServiceImpl) UpdateUserByAccountId(accountId uint64, user domains.User) (uint64, error) {
	//isAccountExist, _ := s.Service.AccountService.IsExistById(accountId)
	//if !isAccountExist {
	//	return 0, errors.New("Не удалось обновить данные пользователя")
	//}
	//var userByAccount *domains.User
	//userByAccount, err := s.Service.UserService.GetUserByAccountId(accountId)
	//if err != nil {
	//	s.Service.Logger.Error(err)
	//	return 0, errors.New("Не удалось обновить данные пользователя")
	//}
	//if userByAccount == nil {
	//	if err != nil {
	//		s.Service.Logger.Error(err)
	//		return 0, errors.New("Не удалось обновить данные пользователя")
	//	}
	//}
	//userByAccount.Name = user.Name
	//userByAccount.Surname = user.Surname
	//userByAccount.Patronymic = user.Patronymic
	//
	//if err = s.Service.UserService.Update(userByAccount); err != nil {
	//	s.Service.Logger.Error(err)
	//	return 0, errors.New("Не удалось обновить данные пользователя")
	//}
	//
	//return *userByAccount.Id, nil
	return 0, errors.New("Переделать")
}

func (s *AccountServiceImpl) GetByEmail(email string) (*domains.Account, error) {
	return s.repository.AccountRepository.GetByEmail(email)
}
