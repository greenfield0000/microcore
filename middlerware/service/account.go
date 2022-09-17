package service

import (
	"microcore/middlerware/domains"
)

type AccountService interface {
	Create(account domains.Account) (id uint64, err error)
	GetAccountInfo(id uint64) (*domains.Account, error)
	All() ([]*domains.Account, error)
	IsExistById(id uint64) (bool, error)
	GetByEmail(email string) (*domains.Account, error)
	DeleteById(id uint64) (bool, error)
	AddAchievement(accountId uint64, achievementId uint64) (id uint64, err error)
	UpdateUserByAccountId(accountId uint64, user domains.User) (uint64, error)
}
