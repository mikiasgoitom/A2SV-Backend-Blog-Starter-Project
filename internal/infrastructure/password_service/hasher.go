package passwordservice

import (
	"blog/internal/domain/contract"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct{}

// check if IHasher was implemented at compile time
var _ contract.IHasher = (*Hasher)(nil)

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (h *Hasher) ComparePasswordHash(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return fmt.Errorf("password verification failed\n")
		}
		return fmt.Errorf("failed to check password hash: %w\n", err)
	}
	return nil
}
