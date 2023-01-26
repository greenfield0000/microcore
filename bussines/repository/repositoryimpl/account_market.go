package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
	"time"
)

type AccountMarketRepositoryImpl struct {
	db *sqlx.DB
}

func NewAccountMarketRepository(db *sqlx.DB) *AccountMarketRepositoryImpl {
	return &AccountMarketRepositoryImpl{db: db}
}

func (a AccountMarketRepositoryImpl) Create(accountMarket domains.AccountMarket) (uint64, error) {
	row := a.db.QueryRow("insert into account_market (marketid, accountid, stateid) values ($1,$2, $3) returning id",
		accountMarket.MarketId,
		accountMarket.AccountId,
		12,
	)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (a AccountMarketRepositoryImpl) DeleteById(id uint64) (bool, error) {
	_, err := a.db.Exec("delete from account_market where id = $1", id)
	return err == nil, err
}

func (a AccountMarketRepositoryImpl) IsMarketInAccount(marketId uint64, accountId uint64) (bool, error) {
	var count uint64
	if err := a.db.Get(&count,
		`select count(*)
			   from account_market am
			            left join state s on am.stateid = s.id
			   where am.marketid = $1
			     and am.accountid = $2
                 and s.sysname = 'ACQUIRED';
		`,
		marketId,
		accountId,
	); err != nil {
		return false, err
	}
	return count != 0, nil
}

func (a AccountMarketRepositoryImpl) GetMarketsByAccountId(accountId uint64) ([]*domains.Market, error) {
	rows, err := a.db.Query(`
			select m.id, m.sysname, m.name, m.cost, m.level, am.accountid
			from market m
         		left join account_market am on am.marketid = m.id and accountid = $1
			union
			select m.id, m.sysname, m.name, m.cost, m.level, am.accountid
			from market m
         		left join account_market am on am.marketid = m.id
			where am.accountid = $1
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
			Level:   level,
			Cost:    cost,
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
