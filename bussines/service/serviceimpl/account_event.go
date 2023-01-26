package serviceimpl

import (
	"context"
	"errors"
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type AccountEventServiceImpl struct {
	logger     *logrus.Logger
	repository *repository.Repository
}

func NewAccountEventService(repository *repository.Repository, logger *logrus.Logger) service.AccountEventService {
	return &AccountEventServiceImpl{logger: logger, repository: repository}
}

func (ae AccountEventServiceImpl) GetEventsByAccountId(accountId uint64) ([]*domains.Event, error) {
	events, err := ae.repository.AccountEventRepository.GetEventsByAccountId(accountId)
	if err != nil {
		ae.logger.Error(err)
		return nil, errors.New("не удалось получить события для аккаунта")
	}
	return events, err
}

func (ae AccountEventServiceImpl) Create(event domains.AccountEvent) (uint64, error) {
	isExist, err := ae.IsEventInAccount(event.EventId, event.AccountId)
	if isExist {
		return 0, errors.New("данное событие уже зарегистрированно для данного пользователя")
	}
	if err != nil {
		ae.logger.Error(err)
		return 0, errors.New("не удалось зарегистрировать событие")
	}
	createdId, err := ae.repository.AccountEventRepository.Create(event.AccountId, event.EventId)
	if err != nil {
		ae.logger.Error(err)
		return 0, errors.New("не удалось зарегистрировать событие")
	}
	return createdId, err
}

func (ae AccountEventServiceImpl) IsEventInAccount(eventId uint64, accountId uint64) (bool, error) {
	id, err := ae.repository.AccountEventRepository.IsExistByEventIdAndAccountId(eventId, accountId)
	if err != nil {
		ae.logger.Error(err)
		return false, errors.New("Не удалось получить информацию о событии")
	}
	return id, err
}

func (ae AccountEventServiceImpl) Remove(ctx context.Context, event domains.AccountEvent) (bool, error) {
	remove, err := ae.repository.AccountEventRepository.Remove(ctx, event)
	if err != nil {
		ae.logger.Error(err)
		return false, errors.New("При удалении события произошла ошибка")
	}
	return remove, err
}

func (ae AccountEventServiceImpl) GetEvent(eventId uint64, accountId uint64) (*domains.Event, error) {
	event, err := ae.repository.AccountEventRepository.GetEvent(eventId, accountId)
	if err != nil {
		ae.logger.Error(err)
		return nil, errors.New("При получении информации о событии произошла ошибка")
	}
	return event, nil
}
