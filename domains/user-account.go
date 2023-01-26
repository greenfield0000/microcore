package domains

type UserAccount struct {
	Id        *uint64 `json:"id" db:"id"`
	UserId    *uint64 `json:"userId" db:"userid"`
	AccountId *uint64 `json:"accountId" db:"accountid"`
}
