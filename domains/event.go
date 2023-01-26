package domains

import "time"

type Event struct {
	Id      *uint64 `json:"id" db:"id"`
	Name    *string `json:"name" db:"name"`
	Sysname *string `json:"sysname" db:"sysname"`
	// вынести в accountEvent
	State      *State     `json:"state"`
	StartDate  *time.Time `json:"startdate"`
	FinishDate *time.Time `json:"finishdate"`
	Cost       float32    `json:"unit"`
}
