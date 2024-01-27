package domains

import "time"

type Team struct {
	Id       *uint64    `json:"id" db:"id"`
	Name     *string    `json:"name" db:"name"`
	Sysname  *string    `json:"sysname" db:"sysname"`
	CreateAt *time.Time `json:"createAt" db:"createat"`
}
