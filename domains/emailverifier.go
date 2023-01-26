package domains

import "time"

type EmailVerifier struct {
	Id           uint64    `db:"id"`
	Code         string    `db:"code"`
	Email        string    `db:"email"`
	VerifyCodeTo time.Time `db:"verify_code_to"`
	StateId      uint64    `db:"stateid"`
}
