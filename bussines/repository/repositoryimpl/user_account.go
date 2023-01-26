package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
)

type UserAccountRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserAccountRepository(db *sqlx.DB) *UserAccountRepositoryImpl {
	return &UserAccountRepositoryImpl{db: db}
}

func (ua UserAccountRepositoryImpl) Create(userAccount domains.UserAccount) (uint64, error) {
	row := ua.db.QueryRow("insert into user_account (userid, accountid) values ($1,$2) returning id",
		userAccount.UserId,
		userAccount.AccountId,
	)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (ua UserAccountRepositoryImpl) DeleteById(id uint64) (bool, error) {
	_, err := ua.db.Exec("delete from user_account where id = $1", id)
	return err == nil, err
}

func (ua UserAccountRepositoryImpl) GetByAccountId(accountId uint64) (*domains.UserAccount, error) {
	var userAccount domains.UserAccount
	if err := ua.db.Get(&userAccount, "select * from user_account where accountid = $1", accountId); err != nil {
		return nil, err
	}
	return &userAccount, nil
}
