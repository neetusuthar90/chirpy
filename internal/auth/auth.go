package auth

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(hashpassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password))
}
