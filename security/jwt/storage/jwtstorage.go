package storage

type JwtStorage interface {
	IsExistExpiredToken(token string) (bool, error)
	PutExpiredToken(accessToken string, refreshToken string) error
}
