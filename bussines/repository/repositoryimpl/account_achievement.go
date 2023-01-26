package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
)

type AccountAchievementRepositoryImpl struct {
	db *sqlx.DB
}

func NewAccountAchievementRepository(db *sqlx.DB) *AccountAchievementRepositoryImpl {
	return &AccountAchievementRepositoryImpl{db: db}
}

func (aa AccountAchievementRepositoryImpl) Create(accountAchievement domains.AccountAchievement) (uint64, error) {
	row := aa.db.QueryRow("insert into account_achievement (achievementid, accountid) values ($1,$2) returning id",
		accountAchievement.AchievementId,
		accountAchievement.AccountId,
	)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
func (aa AccountAchievementRepositoryImpl) DeleteById(id uint64) (bool, error) {
	_, err := aa.db.Exec("delete from account_achievement where id = $1", id)
	return err == nil, err
}

func (aa AccountAchievementRepositoryImpl) IsAchievementInAccount(achievementId uint64, accountId uint64) (bool, error) {
	var count uint64
	if err := aa.db.Get(&count,
		"SELECT count(*) from account_achievement where achievementid = $1 and accountid = $2",
		achievementId,
		accountId,
	); err != nil {
		return false, err
	}
	return count != 0, nil
}
