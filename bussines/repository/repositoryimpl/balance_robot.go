package repositoryimpl

import (
	"github.com/greenfield0000/microcore/domains"
	"github.com/jmoiron/sqlx"
)

type BalanceRobotRepositoryImpl struct {
	db *sqlx.DB
}

func (b BalanceRobotRepositoryImpl) Create(balanceRobot *domains.BalanceRobot) error {
	row := b.db.QueryRow("insert into balance_robot ( stateid, unit, entity_name, entity_id, balance_id, discriminator) values ( 8, $1, $2, $3, $4, $5) returning id;",
		balanceRobot.Unit,
		balanceRobot.EntityName,
		balanceRobot.EntityId,
		balanceRobot.BalanceId,
		balanceRobot.Discriminator,
	)
	var id uint64
	return row.Scan(&id)
}

func NewBalanceRobotRepository(db *sqlx.DB) *BalanceRobotRepositoryImpl {
	return &BalanceRobotRepositoryImpl{
		db: db,
	}
}
