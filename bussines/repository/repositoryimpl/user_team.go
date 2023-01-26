package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
)

type UserTeamRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserTeamRepository(db *sqlx.DB) *UserTeamRepositoryImpl {
	return &UserTeamRepositoryImpl{db: db}
}

func (u *UserTeamRepositoryImpl) Create(userTeam domains.UserTeam) (uint64, error) {
	row := u.db.QueryRow("insert into user_team (teamid, userid) values ($1,$2) returning id",
		userTeam.TeamId,
		userTeam.UserId,
	)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (u *UserTeamRepositoryImpl) IsUserInTeam(teamId uint64, userId uint64) (bool, error) {
	var count uint64
	if err := u.db.Get(&count,
		"SELECT count(*) from user_team where teamid = $1 and userid = $2",
		teamId,
		userId,
	); err != nil {
		return false, err
	}
	return count != 0, nil
}

func (u *UserTeamRepositoryImpl) UpdateUserTeam(userId uint64, teamId uint64) (uint64, error) {
	row := u.db.QueryRow(`update user_team set teamid = $1 where userid = $2 returning id;`,
		teamId,
		userId,
	)
	var updatedId uint64
	if err := row.Scan(&updatedId); err != nil {
		return 0, err
	}
	return updatedId, nil
}

func (u *UserTeamRepositoryImpl) GetByUserId(accountId uint64) (*domains.UserTeam, error) {
	var userTeam domains.UserTeam
	rows, err := u.db.Queryx("select * from user_team where userid = $1", accountId)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	if err := rows.StructScan(&userTeam); err != nil {
		return nil, err
	}
	return &userTeam, nil
}
