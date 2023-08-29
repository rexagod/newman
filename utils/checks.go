package utils

import (
	"os"
	"log"
	"strings"
)

func CheckIfDev(userId string) bool {
	devIDs := strings.Split(os.Getenv("DEVS"), ",")
	isDev := false
	for _, devID := range devIDs {
		log.Println(devID)
		if devID == userId {
			isDev = true
			break
		}
	}

	return isDev
}
