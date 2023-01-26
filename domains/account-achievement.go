package domains

type AccountAchievement struct {
	Id            *uint64 `json:"id,omitempty" db:"id"`
	AccountId     *uint64 `json:"accountId,omitempty" db:"accountid"`
	AchievementId *uint64 `json:"achievementId,omitempty" db:"achievementid"`
}
