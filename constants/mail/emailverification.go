package constant

import "time"

type EmailVerificationState uint64

const (
	// EmailVerificationStateIdWaiting новое подтверждение (в ожидании)
	EmailVerificationStateIdWaiting EmailVerificationState = iota + 16
	// EmailVerificationStateIdConfirmed подтверждено
	EmailVerificationStateIdConfirmed
	// EmailVerificationStateIdError любая ошибка
	EmailVerificationStateIdError
)

// EmailVerificationLag Лаг времени жизни относительно начала регистрации
const EmailVerificationLag = 15 * time.Minute
