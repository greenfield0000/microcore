package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
)

type BalanceRepositoryImpl struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) *BalanceRepositoryImpl {
	return &BalanceRepositoryImpl{db: db}
}

func (b BalanceRepositoryImpl) Create(balance *domains.Balance) error {
	row := b.db.QueryRow("insert into balance (accountid, unit) values ($1, $2) returning id;", balance.AccountId, balance.Unit)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil
}

func (b BalanceRepositoryImpl) GetBalance(accountId uint64) (*domains.Balance, error) {
	row := b.db.QueryRowx("select b.id, b.accountid, b.unit from balance b where b.accountid = $1;", accountId)
	var balance domains.Balance
	if err := row.StructScan(&balance); err != nil {
		return nil, err
	}
	return &balance, nil
}
