package domains

type AccountEvent struct {
	Id        uint64 `json:"id,omitempty" db:"id"`
	EventId   uint64 `json:"eventId,omitempty" db:"eventId"`
	AccountId uint64 `json:"accountId,omitempty" db:"accountId"`
}
