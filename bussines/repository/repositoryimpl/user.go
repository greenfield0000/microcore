package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryIml struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryIml {
	return &UserRepositoryIml{
		db: db,
	}
}

func (u *UserRepositoryIml) Create(user domains.User) (domains.User, error) {
	row := u.db.QueryRow(
		"INSERT INTO \"user\" (name, surname, patronymic) values ($1, $2, $3) returning id",
		user.Name,
		user.Surname,
		user.Patronymic,
	)
	var id *uint64
	if err := row.Scan(&id); err != nil {
		return user, err
	}
	user.Id = id
	return user, nil
}

func (u *UserRepositoryIml) Update(user *domains.User) error {
	_, err := u.db.Exec(
		"update \"user\" set name = $1, surname = $2, patronymic = $3 where id = $4;",
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepositoryIml) GetById(id uint64) (*domains.User, error) {
	var user domains.User
	rows, err := u.db.Queryx("SELECT * from \"user\" where id = $1", id)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err := err; err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	if err = rows.StructScan(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepositoryIml) All() ([]*domains.User, error) {
	rows, err := u.db.Query("SELECT * from \"user\"")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return []*domains.User{}, err
	}
	var uList []*domains.User
	for rows.Next() {
		var (
			id         uint64
			name       *string
			surname    *string
			patronymic *string
		)
		err := rows.Scan(&id, &name, &surname, &patronymic)
		if err != nil {
			return nil, err
		}
		t := &domains.User{
			Name:       name,
			Surname:    surname,
			Patronymic: patronymic,
		}
		uList = append(uList, t)
	}

	return uList, nil
}

func (u *UserRepositoryIml) IsExistById(id *uint64) (bool, error) {
	var count uint64
	if err := u.db.Get(&count, "SELECT count(*) from \"user\" where id = $1", id); err != nil {
		return false, err
	}
	return count != 0, nil
}

func (u *UserRepositoryIml) DeleteById(id uint64) (bool, error) {
	_, err := u.db.Exec("delete from \"user\" where id = $1", id)
	return err == nil, err
}

func (u *UserRepositoryIml) GetUserByAccountId(accountId uint64) (*domains.User, error) {
	var user domains.User
	if err := u.db.Get(&user, "select u.id, u.name, u.surname, u.patronymic from user_account ua left join \"user\" u on ua.userid = u.id where ua.accountid = $1", accountId); err != nil {
		return nil, err
	}
	return &user, nil
}
