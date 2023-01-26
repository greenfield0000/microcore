package domains

type User struct {
	Id         *uint64 `json:"id" db:"id"`
	Name       *string `json:"name" db:"name"`
	Surname    *string `json:"surname" db:"surname"`
	Patronymic *string `json:"patronymic" db:"patronymic"`
}
