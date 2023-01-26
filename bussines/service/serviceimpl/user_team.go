package serviceimpl

import (
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type UserTeamServiceImpl struct {
	repository *repository.Repository
	logger     *logrus.Logger
}

func NewUserTeamService(repository *repository.Repository, logger *logrus.Logger) service.UserTeamService {
	return &UserTeamServiceImpl{repository: repository, logger: logger}
}

func (u UserTeamServiceImpl) Create(userTeam domains.UserTeam) (uint64, error) {
	return u.repository.UserTeamRepository.Create(userTeam)
}

func (u UserTeamServiceImpl) IsUserInTeam(teamId uint64, userId uint64) (bool, error) {
	return u.repository.UserTeamRepository.IsUserInTeam(teamId, userId)
}

func (u UserTeamServiceImpl) GetUserTeamByUserId(id uint64) (*domains.UserTeam, error) {
	return u.repository.UserTeamRepository.GetByUserId(id)
}

func (u UserTeamServiceImpl) UpdateUserTeam(userId uint64, teamId uint64) (uint64, error) {
	return u.repository.UserTeamRepository.UpdateUserTeam(userId, teamId)
}
