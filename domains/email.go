package domains

import "time"

type Email struct {
	Id         uint64    `db:"id"`
	Email      string    `db:"email"`
	CreateDate time.Time `db:"createdate"`
	StateId    uint64    `db:"stateid"`
	TypeId     uint64    `db:"typeid"`
}
