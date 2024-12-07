package utils

import (
	"github.com/google/uuid"
)

var GlobalNucleiInstanceCount int

// GenerateUUID generates a new UUID
func GenerateUUID() (string, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}

func IncreaseNucleiInstanceCount() {
	GlobalNucleiInstanceCount++
}

func DecreaseNucleiInstanceCount() {
	GlobalNucleiInstanceCount--
}

func GetNucleiInstanceCount() int {
	return GlobalNucleiInstanceCount
}
