package domains

import "time"

type EmailVerifier struct {
	Id           uint64    `db:"id"`
	Email        string    `db:"email"`
	CreateDate   time.Time `db:"createdate"`
	StateId      uint64    `db:"stateid"`
	TypeId       uint64    `db:"typeid"`
	VerifyCodeTo time.Time `db:"verify_code_to"`
}
