package constant

type EmailState uint64

const (
	// EmailStateIdWaiting новое подтверждение (в ожидании)
	EmailStateIdWaiting EmailState = iota + 19
	// EmailStateIdSent отправлено
	EmailStateIdSent
	// EmailStateIdError ошибка
	EmailStateIdError
)
