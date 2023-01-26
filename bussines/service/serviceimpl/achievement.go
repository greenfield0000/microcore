package serviceimpl

import (
	"errors"
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type AchievementServiceImpl struct {
	repository *repository.Repository
}

func NewAchievementService(repository *repository.Repository, logger *logrus.Logger) service.AchievementService {
	return &AchievementServiceImpl{
		repository: repository,
	}
}

func (a AchievementServiceImpl) GetAchievementInfo(id uint64) (*domains.Achievement, error) {
	return a.repository.AchievementRepository.GetById(id)
}

func (a AchievementServiceImpl) Create(achievement domains.Achievement) (uint64, error) {
	//exist, err := a.IsExistBySysName(*achievement.Sysname)
	//if err != nil {
	//	a.Service.Logger.Error(err)
	//	return 0, err
	//}
	//if exist {
	//	return 0, errors.New("Achievement already exist")
	//}
	//return a.repository.AchievementRepository.Create(achievement)
	return 0, errors.New("Переделать")
}

func (a AchievementServiceImpl) All() ([]*domains.Achievement, error) {
	return a.repository.AchievementRepository.All()
}

func (a AchievementServiceImpl) IsExistBySysName(sysName string) (bool, error) {
	return a.repository.AchievementRepository.IsExistBySysName(sysName)
}
