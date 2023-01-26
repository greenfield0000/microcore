package repositoryimpl

import (
	"context"
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
	"time"
)

type EventRepositoryImpl struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepositoryImpl {
	return &EventRepositoryImpl{db: db}
}

func (e EventRepositoryImpl) All() ([]*domains.Event, error) {
	rows, err := e.db.Query(`
			select e.id         as id,
    		   e.name       as name,
    		   e.sysname    as sysname,
    		   e.startdate  as startdate,
    		   e.finishdate as finishdate
			from event e
			order by id
	`)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return []*domains.Event{}, err
	}
	var aList []*domains.Event
	for rows.Next() {
		var (
			id         *uint64
			name       *string
			sysname    *string
			startdate  *time.Time
			finishdate *time.Time
		)
		err := rows.Scan(
			&id,
			&name,
			&sysname,
			&startdate,
			&finishdate,
		)
		if err != nil {
			return nil, err
		}
		t := &domains.Event{
			Id:         id,
			Name:       name,
			Sysname:    sysname,
			StartDate:  startdate,
			FinishDate: finishdate,
		}
		aList = append(aList, t)
	}
	return aList, nil
}

func (e EventRepositoryImpl) Create(ctx context.Context, event domains.Event) error {
	row := e.db.QueryRow("insert into event (name, sysname, startdate, finishdate, cost) values ($1, $2, $3, $4, $5) returning id;",
		event.Name,
		event.Sysname,
		event.StartDate,
		event.FinishDate,
		event.Cost,
	)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil
}

func (e EventRepositoryImpl) FindById(ctx context.Context, eventId uint64) (*domains.Event, error) {
	rowX := e.db.QueryRowx("select * from event e where e.id = $1;", eventId)
	var event domains.Event
	err := rowX.StructScan(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e EventRepositoryImpl) Update(ctx context.Context, event domains.Event) error {
	_, err := e.db.Exec("update event set name = $1, sysname = $2, startdate = $3, finishdate = $4, cost = $5 where id = $6;",
		event.Name,
		event.Sysname,
		event.StartDate,
		event.FinishDate,
		event.Cost,
		event.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (e EventRepositoryImpl) DeleteById(ctx context.Context, eventId uint64) error {
	_, err := e.db.Exec("delete from event where id = $1", eventId)
	return err
}
