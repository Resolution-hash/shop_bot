package utils

import "github.com/google/uuid"

func GenereteUniqueFileName(extension string) string {
	return uuid.New().String() + extension
}
