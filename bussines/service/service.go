package service

import (
	"context"

	"github.com/greenfield0000/microcore/domains"
	"github.com/sirupsen/logrus"
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

type UserService interface {
	GetUserInfo(id uint64) (*domains.User, error)
	Create(user domains.User) (domains.User, error)
	All() ([]*domains.User, error)
	IsExistById(id uint64) (bool, error)
	DeleteById(id uint64) (bool, error)
	GetUserByAccountId(accountId uint64) (*domains.User, error)
	Update(user *domains.User) error
}

type TeamService interface {
	GetTeamInfo(id uint64) (*domains.Team, error)
	Create(team domains.Team) (uint64, error)
	All() ([]*domains.Team, error)
	IsExistBySysName(sysName string) (bool, error)
	IsExistById(id uint64) (bool, error)
	AddUser(teamId uint64, userId uint64) (uint64, error)
}

type AchievementService interface {
	GetAchievementInfo(id uint64) (*domains.Achievement, error)
	Create(team domains.Achievement) (uint64, error)
	All() ([]*domains.Achievement, error)
	IsExistBySysName(sysName string) (bool, error)
}

type MarketService interface {
	GetMarketInfo(id uint64) (*domains.Market, error)
	Create(team domains.Market) (uint64, error)
	All() ([]*domains.Market, error)
	IsExistBySysName(sysName string) (bool, error)
	IsExistById(id uint64) (bool, error)
	GetMarketsByAccountId(accountId uint64) (map[string][]*domains.Market, error)
}

type UserTeamService interface {
	Create(userTeam domains.UserTeam) (uint64, error)
	IsUserInTeam(teamId uint64, userId uint64) (bool, error)
	GetUserTeamByUserId(id uint64) (*domains.UserTeam, error)
	UpdateUserTeam(userId uint64, teamId uint64) (uint64, error)
}

type UserAccountService interface {
	Create(userAccount domains.UserAccount) (uint64, error)
	DeleteById(id uint64) (bool, error)
	GetUserAccountByAccountId(accountId uint64) (*domains.UserAccount, error)
}

type AccountAchievementService interface {
	Create(userAccount domains.AccountAchievement) (uint64, error)
	DeleteById(id uint64) (bool, error)
	IsAchievementInAccount(achievementId uint64, accountId uint64) (bool, error)
}

type AccountMarketService interface {
	Create(userAccount domains.AccountMarket) (uint64, error)
	DeleteById(id uint64) (bool, error)
	IsMarketInAccount(marketId uint64, accountId uint64) (bool, error)
}

type AccountEventService interface {
	GetEventsByAccountId(accountId uint64) ([]*domains.Event, error)
	Create(event domains.AccountEvent) (uint64, error)
	IsEventInAccount(eventId uint64, accountId uint64) (bool, error)
	Remove(ctx context.Context, event domains.AccountEvent) (bool, error)
	GetEvent(eventId uint64, accountId uint64) (*domains.Event, error)
}

type EventService interface {
	All() ([]*domains.Event, error)
	Approve(eventId uint64, accountId uint64) error
	Create(ctx context.Context, param domains.Event) error
	Info(ctx context.Context, eventId uint64) (*domains.Event, error)
	Update(ctx context.Context, event domains.Event) error
	Delete(ctx context.Context, eventId uint64) error
}

type BalanceService interface {
	GetBalance(accountId uint64) (*domains.Balance, error)
	Create(balance *domains.Balance) error
}

type BalanceRobotService interface {
	Create(balance *domains.BalanceRobot) error
}

type EmailVerifierService interface {
	CreateCode(email string) (string, error)
	VerifyCode(code string) error
	IsVerifyByEmail(email string) (bool, error)
}

type TemplateService interface {
	Render(path string, data interface{}) (string, error)
}

type MailService interface {
	Send(subject string, to string, message string) error
}

type Service struct {
	logger *logrus.Logger

	AccountService            AccountService
	UserService               UserService
	TeamService               TeamService
	AchievementService        AchievementService
	MarketService             MarketService
	UserTeamService           UserTeamService
	UserAccountService        UserAccountService
	AccountAchievementService AccountAchievementService
	AccountMarketService      AccountMarketService
	AccountEventService       AccountEventService
	EventService              EventService
	BalanceService            BalanceService
	BalanceRobotService       BalanceRobotService
	EmailVerifierService      EmailVerifierService
	TemplateService           TemplateService
	MailService               MailService
}
