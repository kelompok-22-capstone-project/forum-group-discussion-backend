package generator

import "golang.org/x/crypto/bcrypt"

type PasswordGenerator interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

type bcryptPasswordGenerator struct{}

func NewBcryptPasswordGenerator() *bcryptPasswordGenerator {
	return &bcryptPasswordGenerator{}
}

func (b *bcryptPasswordGenerator) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (b *bcryptPasswordGenerator) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
