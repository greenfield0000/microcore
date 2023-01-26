package serviceimpl

import (
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type AccountAchievementServiceImpl struct {
	repository *repository.Repository
}

func NewAccountAchievementService(repository *repository.Repository, logger *logrus.Logger) service.AccountAchievementService {
	return &AccountAchievementServiceImpl{repository: repository}
}

func (aa AccountAchievementServiceImpl) Create(accountAchievement domains.AccountAchievement) (uint64, error) {
	return aa.repository.AccountAchievementRepository.Create(accountAchievement)
}

func (aa AccountAchievementServiceImpl) DeleteById(id uint64) (bool, error) {
	return aa.repository.AccountAchievementRepository.DeleteById(id)
}

func (aa AccountAchievementServiceImpl) IsAchievementInAccount(achievementId uint64, accountId uint64) (bool, error) {
	return aa.repository.AccountAchievementRepository.IsAchievementInAccount(achievementId, accountId)
}
