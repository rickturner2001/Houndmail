package utils

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

func LogAndError(message string) error {
	log.Printf("%s", message)
	return fmt.Errorf("%s", message)
}

func GenToken(size int, symbols bool) string {
	var chars string
	if symbols {
		chars = strings.Join([]string{chars, string("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")}, "")
	}
	chars = strings.Join([]string{chars, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"}, "")

	token := make([]byte, size)
	for i := range token {
		token[i] = chars[rand.Intn(len(chars))]
	}
	return string(token)
}
