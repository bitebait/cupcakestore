package utils

import (
	"crypto/rand"
)

type Randomizer struct {
	charset string
}

func NewRandomizer() *Randomizer {
	return &Randomizer{
		charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}
}

func (r *Randomizer) GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = r.charset[b%byte(len(r.charset))]
	}

	return string(bytes), nil
}
