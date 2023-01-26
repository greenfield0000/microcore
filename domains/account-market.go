package domains

type AccountMarket struct {
	Id        *uint64 `json:"id" db:"id"`
	AccountId *uint64 `json:"accountId" db:"accountid"`
	MarketId  *uint64 `json:"marketId" db:"marketid"`
}
