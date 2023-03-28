package repositoryimpl

import (
	constant "github.com/greenfield0000/microcore/constants/email"
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
	"time"
)

type EmailVerifierRepositoryImpl struct {
	db *sqlx.DB
}

func NewEmailVerifierRepository(db *sqlx.DB) EmailVerifierRepositoryImpl {
	return EmailVerifierRepositoryImpl{
		db: db,
	}
}

func (e EmailVerifierRepositoryImpl) CreateCode(email string, code string, verifyCodeTo time.Time, stateId constant.EmailVerificationState) error {
	_, err := e.db.Exec(
		"INSERT INTO email_verifier (email, code, verify_code_to, stateid) VALUES ($1, $2, $3, $4)",
		email,
		code,
		verifyCodeTo,
		stateId,
	)
	return err
}

func (e EmailVerifierRepositoryImpl) GetCodeWithStatus(code string, stateId constant.EmailVerificationState) (*domains.EmailVerifier, error) {
	var data domains.EmailVerifier
	err := e.db.Get(&data, "select id, code, email, verify_code_to, stateid from email_verifier where code = $1 and stateid = $2", code, stateId)
	if err != nil {
		return nil, err
	}
	return &data, err
}

func (e EmailVerifierRepositoryImpl) GetCode(code string) (*domains.EmailVerifier, error) {
	var data domains.EmailVerifier
	err := e.db.Get(&data, "select id, code, email, verify_code_to, stateid from email_verifier where code = $1", code)
	if err != nil {
		return nil, err
	}
	return &data, err
}

func (e EmailVerifierRepositoryImpl) SetState(code string, stateId constant.EmailVerificationState) error {
	_, err := e.db.Exec("update email_verifier set stateid = $1 where code = $2;", stateId, code)
	return err
}

func (e EmailVerifierRepositoryImpl) IsVerifyByEmail(email string) (bool, error) {
	var count int
	if err := e.db.Get(&count,
		`select count(*)
				from email_verifier e
						 inner join state s on s.id = e.stateid
						 inner join state_entity se on se.id = s.stateentityid
				where se.sysname = 'EMAIL_VERIFICATION' and e.email = $1 and s.sysname = 'CONFIRMED';
			`, email); err != nil {
		return false, err
	}
	return count != 0, nil
}
