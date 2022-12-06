package main

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

func NewDatabase(cfg DatabaseConfig) (*sqlx.DB, error) {
	connect, err := sqlx.Connect(cfg.Drivername,
		fmt.Sprintf(
			"host=%s user=%s  password=%s dbname=%s port=%s sslmode=%s",
			cfg.Host,
			cfg.Username,
			cfg.Password,
			cfg.Dbname,
			cfg.Port,
			cfg.Sslmode,
		),
	)

	if connect == nil {
		return nil, errors.New("no connections")
	}

	connect.SetMaxIdleConns(1)
	connect.SetMaxOpenConns(3)
	connect.SetConnMaxLifetime(3600 * time.Second)

	if err = connect.Ping(); err != nil {

		return nil, err
	}

	//go func() {
	//	for {
	//		// Connected
	//		fmt.Println(fmt.Sprintf("DatabaseConfig count connection: %v", connect.DB.Stats()))
	//		time.Sleep(10 * time.Second)
	//	}
	//}()

	return connect, nil
}
