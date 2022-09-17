package domains

import "time"

// Account учетная запись
type Account struct {
	Id        uint64    `json:"id" db:"id"`
	Login     *string   `json:"login" db:"login"`
	Password  *string   `json:"password" db:"password"`
	Email     *string   `json:"email" db:"email"`
	Phone     *string   `json:"phone" db:"phone"`
	Createdat time.Time `json:"createdat" db:"createat"`
	Blocked   bool      `json:"blocked" db:"blocked"`
	BalanceId uint64    `json:"balanceId"`
}
