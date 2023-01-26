package domains

type UserTeam struct {
	Id     *uint64 `json:"id" db:"id"`
	UserId *uint64 `json:"userId" db:"userid"`
	TeamId *uint64 `json:"teamId" db:"teamid"`
}
