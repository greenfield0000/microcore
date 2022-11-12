package security

import "golang.org/x/crypto/bcrypt"

type PasswordHash struct{}

// CreateHash - function with create hash for password string
func (p *PasswordHash) CreateHash(password string) (string, error) {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytePassword), err
}

// ComparePassword - function with compare hashed password and input password
func (p *PasswordHash) ComparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
