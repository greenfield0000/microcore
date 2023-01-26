package serviceimpl

import (
	"context"
	"errors"
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
)

type EventServiceImpl struct {
	repository *repository.Repository
	logger     *logrus.Logger
}

func NewEventService(repository *repository.Repository, logger *logrus.Logger) service.EventService {
	return &EventServiceImpl{
		logger:     logger,
		repository: repository,
	}
}

func (e EventServiceImpl) All() ([]*domains.Event, error) {
	return e.repository.EventRepository.All()
}

func (e EventServiceImpl) Approve(eventId uint64, accountId uint64) error {
	// TODO апрувить нужно только в определенном состоянии, чтобы не перебивались предыдщие состояния
	if isExist, err := e.repository.AccountEventRepository.IsExistByEventIdAndAccountId(eventId, accountId); !isExist || err != nil {
		e.logger.Error(err)
		return errors.New("данное событие недоступно для данного пользователя")
	}
	event, err := e.repository.AccountEventRepository.GetEvent(eventId, accountId)
	if err != nil {
		e.logger.Error(err)
		return err
	}
	// Получаем баланс аккаунта
	balance, err := e.repository.BalanceRepository.GetBalance(accountId)
	if err != nil {
		e.logger.Error(err)
		return err
	}
	d := "up"
	entityName := "event"
	robot := &domains.BalanceRobot{
		Unit:          &event.Cost,
		EntityName:    &entityName,
		EntityId:      &eventId,
		BalanceId:     &balance.Id,
		Discriminator: &d,
	}
	if err = e.repository.BalanceRobotRepository.Create(robot); err != nil {
		e.logger.Error(err)
		return errors.New("уже находится в обработке")
	}

	if err = e.repository.AccountEventRepository.Approve(eventId, accountId); err != nil {
		e.logger.Error(err)
		return errors.New("При подтверждении произошла ошибка")
	}
	return nil
}

func (e EventServiceImpl) Create(ctx context.Context, event domains.Event) error {
	if err := e.repository.EventRepository.Create(ctx, event); err != nil {
		e.logger.Error(err)
		return errors.New("Не удалось создать событие!")
	}
	return nil
}

func (e EventServiceImpl) Info(ctx context.Context, eventId uint64) (*domains.Event, error) {
	event, err := e.repository.EventRepository.FindById(ctx, eventId)
	if err != nil {
		e.logger.Error(err)
		return nil, errors.New("При получении информации о событии произошла ошибка!")
	}
	return event, nil
}

func (e EventServiceImpl) Update(ctx context.Context, event domains.Event) error {
	if err := e.repository.EventRepository.Update(ctx, event); err != nil {
		e.logger.Error(err)
		return errors.New("Не удалось обновить событие!")
	}
	return nil
}

func (e EventServiceImpl) Delete(ctx context.Context, eventId uint64) error {
	if err := e.repository.EventRepository.DeleteById(ctx, eventId); err != nil {
		return errors.New("Не удалось удалить событие!")
	}
	return nil
}
