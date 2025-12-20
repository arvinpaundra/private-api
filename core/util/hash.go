package util

import "golang.org/x/crypto/bcrypt"

func HashString(str string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func CompareHashAndString(hash, str string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}
