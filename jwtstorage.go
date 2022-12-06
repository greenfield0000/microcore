package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const jwtTableName = "jwtConfig"

type JwtStorage interface {
	IsExistExpiredToken(token string) (bool, error)
	PutExpiredToken(accessToken string, refreshToken string) error
}

type JwtTable struct {
	Id           int64  `db:"id"`
	AccountId    int64  `db:"account_id"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
}

type databaseJwtStorage struct {
	db *sqlx.DB
}

func NewDatabaseJwtStorage(db *sqlx.DB) JwtStorage {
	return &databaseJwtStorage{db: db}
}

func (d databaseJwtStorage) IsExistExpiredToken(token string) (bool, error) {
	var count int64
	if err := d.db.Get(&count, fmt.Sprintf("select count(*) from %s where access_token = $1 or refresh_token = $2", jwtTableName), token, token); err != nil {
		return true, err
	}
	return count != 0, nil
}

func (d databaseJwtStorage) PutExpiredToken(accessToken string, refreshToken string) error {
	query, err := d.db.Query(fmt.Sprintf("INSERT INTO %s (account_id, access_token, refresh_token) VALUES ($1, $2, $3);", jwtTableName), -1, accessToken, refreshToken)
	if query != nil {
		defer query.Close()
	}
	return err
}
