package domains

import "time"

type Market struct {
	Id         *uint64    `json:"id" db:"id"`
	Name       *string    `json:"name" db:"name"`
	Sysname    *string    `json:"sysname" db:"sysname"`
	Level      *int       `json:"level" db:"level"`
	Cost       *float32   `json:"cost" db:"cost"`
	State      *State     `json:"state"`
	CreateDate *time.Time `json:"createdate"`
	UpdateDate *time.Time `json:"updatedate"`
}
