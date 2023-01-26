package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
	"time"
)

type MarketRepositoryImpl struct {
	db *sqlx.DB
}

func NewMarketRepository(db *sqlx.DB) *MarketRepositoryImpl {
	return &MarketRepositoryImpl{db: db}
}

func (m MarketRepositoryImpl) Create(market domains.Market) (uint64, error) {
	row := m.db.QueryRow("insert into market (name ,sysname ,level ,cost) values ($1,$2,$3,$4) returning id",
		market.Name,
		market.Sysname,
		market.Level,
		market.Cost,
	)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (m MarketRepositoryImpl) GetById(id uint64) (*domains.Market, error) {
	var market domains.Market
	if err := m.db.Get(&market, "SELECT * from market where id = $1", id); err != nil {
		return nil, err
	}
	return &market, nil
}

func (m MarketRepositoryImpl) All() ([]*domains.Market, error) {
	rows, err := m.db.Query("SELECT * from market")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return []*domains.Market{}, err
	}
	var mList []*domains.Market
	for rows.Next() {
		var (
			id         *uint64
			name       *string
			sysname    *string
			level      *int
			cost       *float32
			createDate *time.Time
			updateDate *time.Time
		)
		err := rows.Scan(&id, &name, &sysname, &level, &cost, &createDate, &updateDate)
		if err != nil {
			return nil, err
		}
		t := &domains.Market{
			Id:         id,
			Name:       name,
			Sysname:    sysname,
			Level:      level,
			Cost:       cost,
			CreateDate: createDate,
			UpdateDate: updateDate,
		}
		mList = append(mList, t)
	}

	return mList, nil
}

func (m MarketRepositoryImpl) IsExistBySysName(sysName string) (bool, error) {
	var count uint64
	if err := m.db.Get(&count, "SELECT count(*) from market where sysname = $1", sysName); err != nil {
		return false, err
	}
	return count != 0, nil
}

func (m MarketRepositoryImpl) GetMarketsByAccountId(accountId uint64) ([]*domains.Market, error) {
	rows, err := m.db.Query(`
			with market as (select m.id, m.name, m.sysname, m.cost, m.level, am.accountid, s.id, s.name, s.sysname as stateSysName, am.createdate, am.updatedate
			from market m
         		left join account_market am on am.marketid = m.id and accountid = $1
			    left join state s on am.stateid = s.id
			union
			select m.id, m.name, m.sysname, m.cost, m.level, am.accountid, s.id, s.name, s.sysname, am.createdate, am.updatedate
			from market m
         		left join account_market am on am.marketid = m.id
			    left join state s on am.stateid = s.id
			where am.accountid = $1)
            select * from market m
			`, accountId)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return []*domains.Market{}, err
	}
	var mList []*domains.Market
	for rows.Next() {
		var (
			id           *uint64
			name         *string
			sysname      *string
			cost         *float32
			level        *int
			acId         *uint64
			stateId      *uint64
			stateName    *string
			stateSysName *string
			createDate   *time.Time
			updateDate   *time.Time
		)
		err := rows.Scan(&id, &name, &sysname, &cost, &level, &acId, &stateId, &stateName, &stateSysName, &createDate, &updateDate)
		if err != nil {
			return nil, err
		}
		t := &domains.Market{
			Id:      id,
			Name:    name,
			Sysname: sysname,
			Cost:    cost,
			Level:   level,
			State: &domains.State{
				Id:      stateId,
				Name:    stateName,
				Sysname: stateSysName,
			},
			CreateDate: createDate,
			UpdateDate: updateDate,
		}
		mList = append(mList, t)
	}

	return mList, nil
}

func (m MarketRepositoryImpl) IsExistById(id uint64) (bool, error) {
	var count uint64
	if err := m.db.Get(&count, "SELECT count(*) from market where id = $1", id); err != nil {
		return false, err
	}
	return count != 0, nil
}
