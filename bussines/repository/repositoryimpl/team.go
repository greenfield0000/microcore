package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
)

type TeamRepositoryIml struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *TeamRepositoryIml {
	return &TeamRepositoryIml{
		db: db,
	}
}

func (u *TeamRepositoryIml) Create(team domains.Team) (uint64, error) {
	row := u.db.QueryRow(
		"INSERT INTO team (name, sysname) values ($1, $2) returning id",
		team.Name,
		team.Sysname,
	)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (u *TeamRepositoryIml) GetById(id uint64) (*domains.Team, error) {
	var team domains.Team
	if err := u.db.Get(&team, "SELECT * from team where id = $1", id); err != nil {
		return nil, err
	}
	return &team, nil
}

func (u *TeamRepositoryIml) All() ([]*domains.Team, error) {
	rows, err := u.db.Query("SELECT * from team")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return []*domains.Team{}, err
	}
	var tList []*domains.Team
	for rows.Next() {
		var (
			id      *uint64
			name    *string
			sysname *string
		)
		err := rows.Scan(&id, &name, &sysname)
		if err != nil {
			return nil, err
		}
		t := &domains.Team{
			Id:      id,
			Name:    name,
			Sysname: sysname,
		}
		tList = append(tList, t)
	}

	return tList, nil
}

func (u *TeamRepositoryIml) IsExistBySysName(sysName string) (bool, error) {
	var count uint64
	if err := u.db.Get(&count, "SELECT count(*) from team where sysname = $1", sysName); err != nil {
		return false, err
	}
	return count != 0, nil
}

func (u *TeamRepositoryIml) IsExistById(id uint64) (bool, error) {
	var count uint64
	if err := u.db.Get(&count, "SELECT count(*) from team where id = $1", id); err != nil {
		return false, err
	}
	return count != 0, nil
}
