package helpers

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.New("erro ao gerar hash da senha")
	}

	return string(hashedPassword), nil
}

func ParseStringToID(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, strconv.IntSize)

	if err != nil || id == 0 {
		return 0, errors.New("ID inválido")
	}

	return uint(id), nil
}
