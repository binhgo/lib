package util

import (
	"github.com/satori/go.uuid"
)

func GetUUID() (string, error) {
	UUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return UUID.String(), nil
}
