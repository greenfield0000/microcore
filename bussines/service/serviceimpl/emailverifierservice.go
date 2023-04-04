package serviceimpl

import (
	"errors"
	"github.com/google/uuid"
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	constant "github.com/greenfield0000/microcore/constants/email"
	"github.com/sirupsen/logrus"
	"log"

	"time"
)

type EmailVerifierService struct {
	repository *repository.Repository
	logger     *logrus.Logger
}

// NewEmailVerifierService ...
func NewEmailVerifierService(repository *repository.Repository, logger *logrus.Logger) service.EmailVerifierService {
	return &EmailVerifierService{
		repository: repository,
		logger:     logger,
	}
}

// CreateCode ...
func (e EmailVerifierService) CreateCode(email string) (string, error) {
	if ok, _ := e.IsVerifyByEmail(email); !ok {
		code := uuid.New().String()
		err := e.repository.EmailVerifierRepository.CreateCode(email, code, time.Now().Add(constant.EmailVerificationLag), constant.EmailVerificationStateIdWaiting)
		if err != nil {
			return "", errors.New("Не удалось создать код подтверждения")
		}
		return code, nil
	}
	return "", nil
}

// VerifyCode ...
func (e EmailVerifierService) VerifyCode(code string) error {
	// нам нужно проверить, существует ли такой код
	data, err := e.repository.EmailVerifierRepository.GetCode(code)
	if err != nil {
		log.Printf(err.Error())
		return errors.New("Не удалось проверить код")
	}
	if data.StateId == uint64(constant.EmailVerificationStateIdConfirmed) {
		return nil
	}
	if data.StateId != uint64(constant.EmailVerificationStateIdWaiting) {
		return errors.New("Данный код нельзя подтвердить")
	}
	now := time.Now()
	if now.UTC().After(data.VerifyCodeTo.UTC()) {
		err := e.repository.EmailVerifierRepository.SetState(code, constant.EmailVerificationStateIdError)
		if err != nil {
			log.Printf(err.Error())
			return errors.New("Не удалось подтвердить почту")
		}
		return errors.New("Данный код истек")
	}
	err = e.repository.EmailVerifierRepository.SetState(code, constant.EmailVerificationStateIdConfirmed)
	if err != nil {
		log.Printf(err.Error())
		return errors.New("Не удалось подтвердить почту")
	}
	return nil
}

// IsVerifyByEmail ...
func (e EmailVerifierService) IsVerifyByEmail(email string) (bool, error) {
	ok, err := e.repository.EmailVerifierRepository.IsVerifyByEmail(email)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Требуется подтвердить почту")
	}
	return true, nil
}
