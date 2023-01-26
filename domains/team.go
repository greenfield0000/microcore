package domains

type Team struct {
	Id      *uint64 `json:"id" db:"id"`
	Name    *string `json:"name" db:"name"`
	Sysname *string `json:"sysname" db:"sysname"`
}
