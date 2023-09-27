package utils

import (
	"fmt"
	"log"
)

func LogAndError(message string) error {
	log.Printf("%s", message)
	return fmt.Errorf("%s", message)
}
