package random

import (
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]rune, length)
	for i := range result {
		result[i] = rune(chars[rnd.Intn(len(chars))])
	}
	return string(result)
}
