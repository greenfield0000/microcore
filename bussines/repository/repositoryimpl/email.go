package repositoryimpl

import (
	"github.com/greenfield0000/microcore/bussines/repository"
	constant "github.com/greenfield0000/microcore/constants/email"
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
)

type EmailRepositoryImpl struct {
	db *sqlx.DB
}

func NewEmailRepository(db *sqlx.DB) repository.EmailRepository {
	return &EmailRepositoryImpl{
		db: db,
	}
}

func (e EmailRepositoryImpl) GetMailForVerification() ([]domains.Email, error) {
	var data []domains.Email
	err := e.db.Select(&data, `
		select id, email, createdate from email
		where stateid = $1 and typeid = $2
		fetch FIRST 50 ROWS only;
	`, constant.EmailStateIdWaiting, constant.EmailVerificationType)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (e EmailRepositoryImpl) SetState(id uint64, stateId constant.EmailState) error {
	_, err := e.db.Exec("update email set stateid = $1 where id = $2;", stateId, id)
	return err
}
