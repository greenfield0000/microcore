package serviceimpl

import (
	"errors"
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type TeamServiceImpl struct {
	logger *logrus.Logger

	repository *repository.Repository
}

func NewTeamService(repository *repository.Repository, logger *logrus.Logger) service.TeamService {
	return &TeamServiceImpl{logger: logger, repository: repository}
}

func (t *TeamServiceImpl) GetTeamInfo(id uint64) (*domains.Team, error) {
	return t.repository.TeamRepository.GetById(id)
}

func (t *TeamServiceImpl) Create(team domains.Team) (uint64, error) {
	exist, err := t.IsExistBySysName(*team.Sysname)
	if err != nil {
		t.logger.Error(err)
		return 0, err
	}
	if exist {
		return 0, errors.New("Team already exist")
	}
	return t.repository.TeamRepository.Create(team)
}

func (t *TeamServiceImpl) AddUser(teamId uint64, userId uint64) (uint64, error) {
	// Проверка на наличие того и другого
	isTeamExist, _ := t.IsExistById(teamId)
	if !isTeamExist {
		return 0, errors.New("Не удалось добавить пользователя в команду")
	}
	//isUserExist, _ := t.service.UserService.IsExistById(userId)
	//if !isUserExist {
	//	return 0, errors.New("Не удалось добавить пользователя в команду")
	//}
	isUserInTeam, _ := t.repository.IsUserInTeam(teamId, userId)
	if isUserInTeam {
		return 0, errors.New("Пользователь уже состоит в этой команде")
	}
	return t.repository.UserTeamRepository.Create(domains.UserTeam{
		UserId: &userId,
		TeamId: &teamId,
	})
}

func (t *TeamServiceImpl) All() ([]*domains.Team, error) {
	return t.repository.TeamRepository.All()
}

func (t *TeamServiceImpl) IsExistBySysName(sysName string) (bool, error) {
	return t.repository.TeamRepository.IsExistBySysName(sysName)
}

func (t *TeamServiceImpl) IsExistById(id uint64) (bool, error) {
	return t.repository.TeamRepository.IsExistById(id)
}
