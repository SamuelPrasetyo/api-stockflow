package security

import "golang.org/x/crypto/bcrypt"

// PasswordHash handles password hashing operations
type PasswordHash interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type bcryptPasswordHash struct {
	cost int
}

// NewBcryptPasswordHash creates a new bcrypt password hash handler
func NewBcryptPasswordHash() PasswordHash {
	return &bcryptPasswordHash{
		cost: bcrypt.DefaultCost,
	}
}

func (b *bcryptPasswordHash) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (b *bcryptPasswordHash) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
