package repository

import (
	"context"
	constant "github.com/greenfield0000/microcore/constants/email"
	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
	"time"
)

type AccountRepository interface {
	Create(account domains.Account) (uint64, error)
	GetById(id uint64) (*domains.Account, error)
	All() ([]*domains.Account, error)
	IsExist(account uint64) (bool, error)
	DeleteById(id uint64) (bool, error)
	GetByEmail(email string) (*domains.Account, error)
}

type UserRepository interface {
	Create(user domains.User) (domains.User, error)
	GetById(id uint64) (*domains.User, error)
	All() ([]*domains.User, error)
	IsExistById(id *uint64) (bool, error)
	DeleteById(id uint64) (bool, error)
	GetUserByAccountId(accountId uint64) (*domains.User, error)
	Update(user *domains.User) error
}

type TeamRepository interface {
	Create(team domains.Team) (uint64, error)
	GetById(id uint64) (*domains.Team, error)
	All() ([]*domains.Team, error)
	IsExistBySysName(sysName string) (bool, error)
	IsExistById(id uint64) (bool, error)
}

type AchievementRepository interface {
	Create(achievement domains.Achievement) (uint64, error)
	GetById(id uint64) (*domains.Achievement, error)
	All() ([]*domains.Achievement, error)
	IsExistBySysName(sysName string) (bool, error)
	IsExistById(id uint64) (bool, error)
}

type MarketRepository interface {
	Create(market domains.Market) (uint64, error)
	GetById(id uint64) (*domains.Market, error)
	All() ([]*domains.Market, error)
	IsExistBySysName(sysName string) (bool, error)
	IsExistById(id uint64) (bool, error)
	GetMarketsByAccountId(accountId uint64) ([]*domains.Market, error)
}

type UserTeamRepository interface {
	Create(userTeam domains.UserTeam) (uint64, error)
	IsUserInTeam(teamId uint64, userId uint64) (bool, error)
	GetByUserId(id uint64) (*domains.UserTeam, error)
	UpdateUserTeam(userId uint64, teamId uint64) (uint64, error)
}

type UserAccountRepository interface {
	Create(userAccount domains.UserAccount) (uint64, error)
	DeleteById(id uint64) (bool, error)
	GetByAccountId(accountId uint64) (*domains.UserAccount, error)
}

type AccountAchievementRepository interface {
	Create(userAccount domains.AccountAchievement) (uint64, error)
	DeleteById(id uint64) (bool, error)
	IsAchievementInAccount(achievementId uint64, accountId uint64) (bool, error)
}

type AccountMarketRepository interface {
	Create(userAccount domains.AccountMarket) (uint64, error)
	DeleteById(id uint64) (bool, error)
	IsMarketInAccount(marketId uint64, accountId uint64) (bool, error)
	GetMarketsByAccountId(accountId uint64) ([]*domains.Market, error)
}

type AccountEventRepository interface {
	GetEventsByAccountId(accountId uint64) ([]*domains.Event, error)
	Create(accountId uint64, eventId uint64) (uint64, error)
	IsExistByEventIdAndAccountId(eventId uint64, accountId uint64) (bool, error)
	Approve(eventId uint64, accountId uint64) error
	Remove(ctx context.Context, event domains.AccountEvent) (bool, error)
	GetEvent(eventId uint64, accountId uint64) (*domains.Event, error)
}

type EventRepository interface {
	All() ([]*domains.Event, error)
	Create(ctx context.Context, event domains.Event) error
	FindById(ctx context.Context, eventId uint64) (*domains.Event, error)
	Update(ctx context.Context, event domains.Event) error
	DeleteById(ctx context.Context, eventId uint64) error
}

type BalanceRepository interface {
	GetBalance(accountId uint64) (*domains.Balance, error)
	Create(balance *domains.Balance) error
}

type BalanceRobotRepository interface {
	Create(balanceRobot *domains.BalanceRobot) error
}

type EmailVerifierRepository interface {
	CreateCode(email string, code string, verifyCodeTo time.Time, stateId constant.EmailVerificationState) error
	GetCodeWithStatus(code string, stateId constant.EmailVerificationState) (*domains.EmailVerifier, error)
	GetCode(code string) (*domains.EmailVerifier, error)
	SetState(code string, stateId constant.EmailVerificationState) error
}

type EmailRepository interface {
	GetMailForVerification() ([]domains.Email, error)
	SetState(id uint64, stateId constant.EmailState) error
}

type Repository struct {
	logger *logrus.Logger

	AccountRepository
	UserRepository
	TeamRepository
	AchievementRepository
	MarketRepository
	UserTeamRepository
	UserAccountRepository
	AccountAchievementRepository
	AccountMarketRepository
	AccountEventRepository
	EventRepository
	BalanceRepository
	BalanceRobotRepository
	EmailVerifierRepository
	EmailRepository
}
