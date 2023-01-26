package domains

type BalanceRobot struct {
	Id            *uint64  `json:"id"`
	ChannelId     *uint64  `json:"channelId"`
	StateId       *uint64  `json:"stateId"`
	Unit          *float32 `json:"unit"`
	EntityName    *string  `json:"entityName"`
	EntityId      *uint64  `json:"entityId"`
	BalanceId     *uint64  `json:"balanceId"`
	Discriminator *string  `json:"discriminator"`
}
