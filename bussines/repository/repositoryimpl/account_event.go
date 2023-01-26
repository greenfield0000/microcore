package repositoryimpl

import (
	"context"
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
	"time"
)

type AccountEventRepositoryImpl struct {
	db *sqlx.DB
}

func NewAccountEventRepository(db *sqlx.DB) *AccountEventRepositoryImpl {
	return &AccountEventRepositoryImpl{db}
}

func (ae AccountEventRepositoryImpl) GetEvent(eventId uint64, accountId uint64) (*domains.Event, error) {
	rowx := ae.db.QueryRowx(
		`
				select e.id, e.name, e.sysname, s.id, s.name, s.sysname, e.startdate, e.finishdate, e.cost
					from event e
					         left join account_event ae on e.id = ae.eventid
					         left join state s on ae.stateid = s.id
					where accountid = $1 and ae.eventid = $2
				`,
		accountId,
		eventId,
	)
	var (
		id           *uint64
		name         *string
		sysname      *string
		stateId      *uint64
		stateName    *string
		stateSysname *string
		startdate    *time.Time
		finishdate   *time.Time
		unit         float32
	)
	err := rowx.Scan(&id, &name, &sysname, &stateId, &stateName, &stateSysname, &startdate, &finishdate, &unit)
	if err != nil {
		return nil, err
	}
	return &domains.Event{
		Id:      id,
		Name:    name,
		Sysname: sysname,
		State: &domains.State{
			Id:      stateId,
			Name:    stateName,
			Sysname: stateSysname,
		},
		StartDate:  startdate,
		FinishDate: finishdate,
		Cost:       unit,
	}, nil
}

func (aer AccountEventRepositoryImpl) GetEventsByAccountId(accountId uint64) ([]*domains.Event, error) {
	rows, err := aer.db.Queryx(
		`
					select e.id, e.name, e.sysname, e.startdate, e.finishdate, s.id, s.name, s.sysname, ae.accountid
					from event e
					         left join account_event ae on ae.eventid = e.id and accountid = $1
					         left join state s on ae.stateid = s.id
					union
					select e.id, e.name, e.sysname, e.startdate, e.finishdate, s.id, s.name, s.sysname, ae.accountid
					from event e
					         left join account_event ae on ae.eventid = e.id
					         left join state s on ae.stateid = s.id
					where ae.accountid = $1
				`,
		accountId,
	)

	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return []*domains.Event{}, err
	}

	var eventList []*domains.Event
	for rows.Next() {
		var (
			id           *uint64
			name         *string
			sysname      *string
			startdate    *time.Time
			finishdate   *time.Time
			stateId      *uint64
			stateName    *string
			stateSysname *string
			acId         *uint64
		)
		err := rows.Scan(&id, &name, &sysname, &startdate, &finishdate, &stateId, &stateName, &stateSysname, &acId)
		if err != nil {
			return []*domains.Event{}, err
		}
		eventList = append(eventList, &domains.Event{
			Id:      id,
			Name:    name,
			Sysname: sysname,
			State: &domains.State{
				Id:      stateId,
				Name:    stateName,
				Sysname: stateSysname,
			},
			StartDate:  startdate,
			FinishDate: finishdate,
		})
	}
	return eventList, nil
}

func (aer AccountEventRepositoryImpl) Create(accountId uint64, eventId uint64) (uint64, error) {
	row := aer.db.QueryRowx("insert into account_event (eventid, accountid, stateid) VALUES ($1, $2, $3) returning id;",
		eventId,
		accountId,
		4,
	)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (aer AccountEventRepositoryImpl) IsExistByEventIdAndAccountId(eventId uint64, accountId uint64) (bool, error) {
	var count int
	if err := aer.db.Get(&count, "select count(*) from account_event where eventid = $1 and accountid = $2 and stateid = 4", eventId, accountId); err != nil {
		return false, err
	}
	return count != 0, nil
}

func (aer AccountEventRepositoryImpl) Approve(eventId uint64, accountId uint64) error {
	_, err := aer.db.Exec("update account_event set stateid = $1 where eventid = $2 and accountid = $3;",
		// TODO: вынести в сущность состояния
		5,
		eventId,
		accountId,
	)
	return err
}

func (aer AccountEventRepositoryImpl) Remove(ctx context.Context, event domains.AccountEvent) (bool, error) {
	var count *int64
	rowx := aer.db.QueryRow("delete from account_event where eventid = $1 and accountid = $2 and stateid = 4 returning id", event.EventId, event.AccountId)
	err := rowx.Scan(&count)
	return count != nil || err == nil, nil
}
