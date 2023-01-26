package domains

type Balance struct {
	Id        uint64 `json:"id"`
	AccountId uint64 `json:"accountId"`
	Unit      uint64 `json:"unit"`
}
