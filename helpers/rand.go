package helpers

import (
	"crypto/rand"
)

const defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Randomizer struct
type Randomizer struct {
	charset string
}

func NewRandomizer() *Randomizer {
	return &Randomizer{
		charset: defaultCharset,
	}
}

func (r *Randomizer) GenerateString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	for i, b := range randomBytes {
		randomBytes[i] = r.charset[b%byte(len(r.charset))]
	}
	return string(randomBytes), nil
}
