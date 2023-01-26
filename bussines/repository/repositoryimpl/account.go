package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
	"time"
)

type AccountRepositoryImpl struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{
		db: db,
	}
}

func (r AccountRepositoryImpl) Create(account domains.Account) (uint64, error) {
	row := r.db.QueryRow(
		"INSERT INTO account (password, email, createat) VALUES ($1, $2, $3) returning id",
		account.Password,
		account.Email,
		account.Createdat,
	)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r AccountRepositoryImpl) GetById(id uint64) (*domains.Account, error) {
	var account domains.Account
	if err := r.db.Get(&account, "SELECT * FROM account where id = $1", id); err != nil {
		return nil, err
	}
	return &account, nil
}

func (r AccountRepositoryImpl) All() ([]*domains.Account, error) {
	rows, err := r.db.Query("SELECT * from account")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return []*domains.Account{}, err
	}
	var acList []*domains.Account
	for rows.Next() {
		var (
			id        uint64
			login     *string
			password  *string
			email     *string
			phone     *string
			createdat time.Time
			blocked   bool
		)
		err := rows.Scan(&id, &login, &password, &email, &phone, &createdat, &blocked)
		if err != nil {
			return nil, err
		}
		ac := &domains.Account{
			Id:        id,
			Login:     login,
			Email:     email,
			Phone:     phone,
			Createdat: createdat,
			Blocked:   blocked,
		}
		acList = append(acList, ac)
	}

	return acList, nil
}

func (r AccountRepositoryImpl) DeleteById(id uint64) (bool, error) {
	_, err := r.db.Exec("delete from account where id = $1", id)
	return err == nil, err
}

func (r AccountRepositoryImpl) IsExist(accountId uint64) (bool, error) {
	var count int
	if err := r.db.Get(&count, "SELECT count(*) from account where id = $1", accountId); err != nil {
		return false, err
	}
	return count != 0, nil
}

func (r AccountRepositoryImpl) GetByEmail(email string) (*domains.Account, error) {
	var account domains.Account
	row := r.db.QueryRowx("SELECT * FROM account where email = $1", email)
	if row.Err() != nil {
		return nil, row.Err()
	}
	err := row.StructScan(&account)
	if err != nil {
		return nil, nil
	}
	return &account, nil
}
