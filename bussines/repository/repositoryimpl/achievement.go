package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
)

type AchievementRepositoryImpl struct {
	db *sqlx.DB
}

func NewAchievementRepository(db *sqlx.DB) *AchievementRepositoryImpl {
	return &AchievementRepositoryImpl{db: db}
}

func (a AchievementRepositoryImpl) Create(achievement domains.Achievement) (uint64, error) {
	row := a.db.QueryRow("insert into achievement (name ,sysname ,level ,cost) values ($1,$2,$3,$4) returning id",
		achievement.Name,
		achievement.Sysname,
		achievement.Level,
		achievement.Cost,
	)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (a AchievementRepositoryImpl) GetById(id uint64) (*domains.Achievement, error) {
	var ac domains.Achievement
	if err := a.db.Get(&ac, "SELECT * from achievement where id = $1", id); err != nil {
		return nil, err
	}
	return &ac, nil
}

func (a AchievementRepositoryImpl) All() ([]*domains.Achievement, error) {
	rows, err := a.db.Query("SELECT * from achievement")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return []*domains.Achievement{}, err
	}
	var aList []*domains.Achievement
	for rows.Next() {
		var (
			id      uint64
			name    *string
			sysname *string
			level   *int
			cost    *float32
		)
		err := rows.Scan(&id, &name, &sysname, &level, &cost)
		if err != nil {
			return nil, err
		}
		t := &domains.Achievement{
			Name:    name,
			Sysname: sysname,
			Level:   level,
			Cost:    cost,
		}
		aList = append(aList, t)
	}

	return aList, nil
}

func (a AchievementRepositoryImpl) IsExistBySysName(sysName string) (bool, error) {
	var count uint64
	if err := a.db.Get(&count, "SELECT count(*) from achievement where sysName = $1", sysName); err != nil {
		return false, err
	}
	return count != 0, nil
}

func (a AchievementRepositoryImpl) IsExistById(id uint64) (bool, error) {
	var count uint64
	if err := a.db.Get(&count, "SELECT count(*) from achievement where id = $1", id); err != nil {
		return false, err
	}
	return count != 0, nil
}
