package constant

type EmailType uint64

const (
	// EmailStateIdWaiting новое подтверждение (в ожидании)
	EmailVerificationType EmailType = iota + 1
)
