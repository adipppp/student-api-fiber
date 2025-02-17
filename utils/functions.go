package utils

import (
	"errors"
	"os"
)

func ValidateEnv() error {
	_, ok := os.LookupEnv("PORT")
	if !ok {
		return errors.New("PORT environment variable not set")
	}
	_, ok = os.LookupEnv("SSO_UI_URL")
	if !ok {
		return errors.New("SSO_UI_URL environment variable not set")
	}
	return nil
}
