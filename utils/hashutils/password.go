package hashutils

import (
	"golang.org/x/crypto/bcrypt"
)

// EncodePassword encodes the given password value with bcrypt and default cost.
func EncodePassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// MatchesPassword returns a true if matched hashedPassword and raw password, otherwise false.
func MatchesPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
