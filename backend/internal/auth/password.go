package auth

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 10

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() BcryptPasswordHasher {
	return BcryptPasswordHasher{}
}

func (BcryptPasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (BcryptPasswordHasher) Compare(passwordHash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
